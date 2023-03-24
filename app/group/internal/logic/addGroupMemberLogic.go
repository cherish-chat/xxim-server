package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/group/groupmodel"
	"github.com/cherish-chat/xxim-server/app/msg/msgmodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"go.opentelemetry.io/otel/propagation"
	"gorm.io/gorm"
	"time"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddGroupMemberLogic {
	return &AddGroupMemberLogic{
		ctx:    xtrace.NewContext(ctx),
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AddGroupMember 添加群成员
func (l *AddGroupMemberLogic) AddGroupMember(in *pb.AddGroupMemberReq) (*pb.AddGroupMemberResp, error) {
	if len(in.UserIds) == 0 {
		return &pb.AddGroupMemberResp{}, nil
	}
	err := xorm.Transaction(l.svcCtx.Mysql(), func(tx *gorm.DB) error {
		var members []*groupmodel.GroupMember
		for _, uid := range in.UserIds {
			member := &groupmodel.GroupMember{
				GroupId:    in.GroupId,
				UserId:     uid,
				CreateTime: time.Now().UnixMilli(),
				Role:       groupmodel.RoleType_MEMBER,
			}
			members = append(members, member)
		}
		err := tx.Model(&groupmodel.GroupMember{}).CreateInBatches(members, 100).Error
		if err != nil {
			// 判断是不是唯一索引冲突
			if !xorm.DuplicateError(err) {
				l.Errorf("InviteFriendToGroup InsertMany error: %v", err)
				return err
			}
		}
		return nil
	}, func(tx *gorm.DB) error {
		return groupmodel.FlushGroupMemberListCache(l.ctx, l.svcCtx.Redis(), in.GroupId)
	}, func(tx *gorm.DB) error {
		return groupmodel.FlushGroupMemberCache(l.ctx, l.svcCtx.Redis(), in.GroupId, in.UserIds...)
	})
	if err != nil {
		l.Errorf("AddGroupMember error: %v", err)
		return &pb.AddGroupMemberResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	utils.RetryProxy(context.Background(), 12, 1*time.Second, func() error {
		// 删除缓存
		{
			err = groupmodel.FlushGroupsByUserIdCache(l.ctx, l.svcCtx.Redis(), in.UserIds...)
			if err != nil {
				l.Errorf("InviteFriendToGroup FlushGroupsByUserIdCache error: %v", err)
				return err
			}
		}
		_, err := l.svcCtx.MsgService().FlushUsersSubConv(l.ctx, &pb.FlushUsersSubConvReq{UserIds: in.UserIds})
		if err != nil {
			l.Errorf("FlushUsersSubConv failed, err: %v", err)
			return err
		}
		_, err = NewSyncGroupMemberCountLogic(l.ctx, l.svcCtx).SyncGroupMemberCount(&pb.SyncGroupMemberCountReq{
			CommonReq: in.GetCommonReq(),
			GroupId:   in.GroupId,
		})
		if err != nil {
			l.Errorf("SyncGroupMemberCount failed, err: %v", err)
			return err
		}
		return nil
	})

	// 申请人发一条消息 自我介绍
	go xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "SendMyInfo", func(ctx context.Context) {
		// 获取我自己的信息
		userBaseInfoMap := make(map[string]*pb.UserBaseInfo)
		{
			userBaseInfo, err := l.svcCtx.UserService().BatchGetUserBaseInfo(ctx, &pb.BatchGetUserBaseInfoReq{
				CommonReq: in.CommonReq,
				Ids:       in.UserIds,
			})
			if err != nil {
				l.Errorf("SendMsg failed, err: %v", err)
				return
			}
			if len(userBaseInfo.UserBaseInfos) == 0 {
				l.Errorf("SendMsg failed, err: %v", err)
				return
			}
			for _, info := range userBaseInfo.UserBaseInfos {
				userBaseInfoMap[info.Id] = info
			}
		}
		var msgDatas []*pb.MsgData
		for _, id := range in.UserIds {
			sender, ok := userBaseInfoMap[id]
			if !ok {
				continue
			}
			data := msgmodel.CreateTextMsgToGroup(
				sender,
				in.GroupId,
				"大家好，我是"+sender.Nickname+"。", msgmodel.MsgOptions{
					OfflinePush:       false,
					StorageForServer:  true,
					StorageForClient:  true,
					UpdateUnreadCount: true,
					NeedDecrypt:       false,
					UpdateConvMsg:     true,
				},
				nil,
				nil,
			).ToMsgData()
			msgDatas = append(msgDatas, data)
		}
		_, err := l.svcCtx.MsgService().SendMsgListAsync(ctx, &pb.SendMsgListReq{
			MsgDataList:  msgDatas,
			DeliverAfter: nil,
			CommonReq:    &pb.CommonReq{Platform: "system"},
		})
		if err != nil {
			l.Errorf("SendMsg failed, err: %v", err)
		}
	}, propagation.MapCarrier{
		"groupId": in.GroupId,
	})
	return &pb.AddGroupMemberResp{}, nil
}
