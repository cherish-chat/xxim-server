package logic

import (
	"context"
	msgservice "github.com/cherish-chat/xxim-server/app/msg/msgService"
	"github.com/cherish-chat/xxim-server/app/msg/msgmodel"
	"github.com/cherish-chat/xxim-server/app/notice/noticemodel"
	"github.com/cherish-chat/xxim-server/app/relation/relationmodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"gorm.io/gorm"
	"time"

	"github.com/cherish-chat/xxim-server/app/relation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchMakeFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchMakeFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchMakeFriendLogic {
	return &BatchMakeFriendLogic{
		ctx:    xtrace.NewContext(ctx),
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BatchMakeFriendLogic) BatchMakeFriend(in *pb.BatchMakeFriendReq) (*pb.BatchMakeFriendResp, error) {
	if len(in.UserIdBList) == 0 {
		return &pb.BatchMakeFriendResp{}, nil
	}
	// 我的好友总数是否已达上限
	{
		var getFriendCountResp *pb.GetFriendCountResp
		var err error
		xtrace.StartFuncSpan(l.ctx, "GetFriendCount", func(ctx context.Context) {
			getFriendCountResp, err = NewGetFriendCountLogic(ctx, l.svcCtx).GetFriendCount(&pb.GetFriendCountReq{
				CommonReq: in.CommonReq,
			})
		})
		if err != nil {
			l.Errorf("GetFriendCount failed, err: %v", err)
			return &pb.BatchMakeFriendResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
		if int64(getFriendCountResp.Count) >= l.svcCtx.ConfigMgr.FriendMaxCount(l.ctx, in.UserIdA) {
			return &pb.BatchMakeFriendResp{CommonResp: pb.NewToastErrorResp(l.svcCtx.T(in.CommonReq.Language, "好友数量已达上限"))}, nil
		}
	}
	now := time.Now().UnixMilli()
	friends := make([]*relationmodel.Friend, 0)
	for _, bId := range in.UserIdBList {
		friends = append(friends, &relationmodel.Friend{
			UserId:     in.UserIdA,
			FriendId:   bId,
			CreateTime: now,
		}, &relationmodel.Friend{
			UserId:     bId,
			FriendId:   in.UserIdA,
			CreateTime: now,
		})
	}
	{
		// 添加好友
		err := xorm.Transaction(l.svcCtx.Mysql(), func(tx *gorm.DB) error {
			err := tx.Model(&relationmodel.Friend{}).Where("userId = ? AND friendId IN (?)", in.UserIdA, in.UserIdBList).Delete(&relationmodel.Friend{}).Error
			if err != nil {
				l.Errorf("DeleteOne failed, err: %v", err)
				return err
			}
			err = tx.Model(&relationmodel.Friend{}).Where("userId IN (?) AND friendId = ?", in.UserIdBList, in.UserIdA).Delete(&relationmodel.Friend{}).Error
			if err != nil {
				l.Errorf("DeleteOne failed, err: %v", err)
				return err
			}
			err = tx.Model(&relationmodel.Friend{}).CreateInBatches(friends, 100).Error
			if err != nil {
				l.Errorf("InsertOne failed, err: %v", err)
				return err
			}
			return nil
		}, func(tx *gorm.DB) error {
			for _, userId := range append(in.UserIdBList, in.UserIdA) {
				notice := &noticemodel.Notice{
					ConvId: pb.HiddenConvIdCommand(),
					UserId: userId,
					Options: noticemodel.NoticeOption{
						StorageForClient: false,
						UpdateConvNotice: false,
					},
					ContentType: pb.NoticeContentType_SyncFriendList,
					Content: utils.AnyToBytes(pb.NoticeContent_SyncFriendList{
						Comment: "acceptAddFriend",
					}),
					UniqueId: "syncFriendList",
					Title:    "",
					Ext:      nil,
				}
				err := notice.Insert(l.ctx, tx, l.svcCtx.Redis())
				if err != nil {
					l.Errorf("insert notice failed, err: %v", err)
					return err
				}
			}
			return nil
		})
		if err != nil {
			l.Errorf("InsertOne failed, err: %v", err)
			return &pb.BatchMakeFriendResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	{
		// 删除缓存
		err := relationmodel.FlushFriendList(l.ctx, l.svcCtx.Redis(), append(in.UserIdBList, in.UserIdA)...)
		if err != nil {
			l.Errorf("FlushFriendList failed, err: %v", err)
		}
		// 刷新订阅
		utils.RetryProxy(context.Background(), 12, 1*time.Second, func() error {
			_, err := l.svcCtx.MsgService().FlushUsersSubConv(l.ctx, &pb.FlushUsersSubConvReq{UserIds: append(in.UserIdBList, in.UserIdA)})
			if err != nil {
				l.Errorf("FlushUsersSubConv failed, err: %v", err)
				return err
			}
			for _, userId := range append(in.UserIdBList, in.UserIdA) {
				_, err = l.svcCtx.NoticeService().GetUserNoticeData(l.ctx, &pb.GetUserNoticeDataReq{
					UserId: userId,
					ConvId: pb.HiddenConvIdCommand(),
				})
				if err != nil {
					l.Errorf("SendNoticeData failed, err: %v", err)
					return err
				}
			}
			return err
		})
		// 接受者发送消息：我们已经是好友了，快来聊天吧
		l.sendMsg(in)
	}
	return &pb.BatchMakeFriendResp{}, nil
}

func (l *BatchMakeFriendLogic) sendMsg(in *pb.BatchMakeFriendReq) {
	go xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "SendMsg", func(ctx context.Context) {
		// 获取接受者info
		var userByIds = make(map[string]*pb.UserBaseInfo)
		batchGetUserBaseInfoResp, err := l.svcCtx.UserService().BatchGetUserBaseInfo(ctx, &pb.BatchGetUserBaseInfoReq{Ids: []string{in.UserIdA}})
		if err != nil {
			l.Errorf("MapUserByIds failed, err: %v", err)
		} else {
			for _, info := range batchGetUserBaseInfoResp.UserBaseInfos {
				userByIds[info.Id] = info
			}
			selfInfo, ok := userByIds[in.UserIdA]
			if ok {
				text := "我们已经是好友了，快来聊天吧"
				if in.SendTextMsgA != nil && *in.SendTextMsgA != "" {
					text = *in.SendTextMsgA
				}
				var msgDatas []*pb.MsgData
				for _, bId := range in.UserIdBList {
					data := msgmodel.CreateTextMsgToUser(
						selfInfo,
						bId,
						l.svcCtx.T(in.CommonReq.Language, text),
						msgmodel.MsgOptions{
							OfflinePush:       true,
							StorageForServer:  true,
							StorageForClient:  true,
							UpdateUnreadCount: true,
							NeedDecrypt:       false,
							UpdateConvMsg:     true,
						},
						&msgmodel.MsgOfflinePush{
							Title:   selfInfo.Nickname,
							Content: text,
							Payload: "",
						},
						nil,
					).ToMsgData()
					msgDatas = append(msgDatas, data)
				}
				_, err = msgservice.SendMsgSync(l.svcCtx.MsgService(), ctx, msgDatas)
				if err != nil {
					l.Errorf("SendMsgSync failed, err: %v", err)
					err = nil
				}
			}
		}
	}, nil)
}
