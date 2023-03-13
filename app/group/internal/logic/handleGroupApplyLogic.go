package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/group/groupmodel"
	msgservice "github.com/cherish-chat/xxim-server/app/msg/msgService"
	"github.com/cherish-chat/xxim-server/app/msg/msgmodel"
	"github.com/cherish-chat/xxim-server/app/notice/noticemodel"
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

type HandleGroupApplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewHandleGroupApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HandleGroupApplyLogic {
	return &HandleGroupApplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// HandleGroupApply 处理群聊申请
func (l *HandleGroupApplyLogic) HandleGroupApply(in *pb.HandleGroupApplyReq) (*pb.HandleGroupApplyResp, error) {
	// 查询 apply
	apply := &groupmodel.GroupApply{}
	err := l.svcCtx.Mysql().Model(&groupmodel.GroupApply{}).Where("id = ?", in.ApplyId).Limit(1).Find(apply).Error
	if err != nil {
		l.Errorf("HandleGroupApply Find error: %v", err)
		return &pb.HandleGroupApplyResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	if apply.Id == "" {
		return &pb.HandleGroupApplyResp{CommonResp: pb.NewSuccessResp()}, nil
	}
	if apply.HandleTime != 0 {
		return &pb.HandleGroupApplyResp{CommonResp: pb.NewSuccessResp()}, nil
	}
	// 获取群里所有的管理员
	groupManagers, err := getAllGroupManager(l.ctx, l.svcCtx, apply.GroupId, true)
	if err != nil {
		l.Errorf("HandleGroupApply getAllGroupManager error: %v", err)
		return &pb.HandleGroupApplyResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	// 判断是否是管理员
	isManager := false
	for _, manager := range groupManagers {
		if manager.UserId == in.CommonReq.UserId {
			isManager = true
			break
		}
	}
	if !isManager {
		return &pb.HandleGroupApplyResp{CommonResp: pb.NewAlertErrorResp(
			l.svcCtx.T(in.CommonReq.Language, "操作失败"),
			l.svcCtx.T(in.CommonReq.Language, "您不是群管理员"),
		)}, nil
	}
	updateMap := map[string]interface{}{
		"result":       in.Result,
		"handleTime":   time.Now().UnixMilli(),
		"handleUserId": in.CommonReq.UserId,
	}
	{
		apply.Result = in.Result
		apply.HandleTime = time.Now().UnixMilli()
		apply.HandleUserId = in.CommonReq.UserId
	}
	err = xorm.Transaction(l.svcCtx.Mysql(), func(tx *gorm.DB) error {
		// 更新 apply
		err := tx.Model(&groupmodel.GroupApply{}).Where("id = ?", in.ApplyId).Updates(updateMap).Error
		if err != nil {
			l.Errorf("HandleGroupApply Update error: %v", err)
			return err
		}
		return nil
	}, func(tx *gorm.DB) error {
		if in.Result == pb.GroupApplyHandleResult_AGREE {
			// 同意 加入群
			// 群成员
			member := &groupmodel.GroupMember{
				GroupId:    apply.GroupId,
				UserId:     apply.UserId,
				CreateTime: time.Now().UnixMilli(),
				Role:       groupmodel.RoleType_MEMBER,
			}
			err := xorm.InsertOne(tx, member)
			if err != nil {
				// 判断是不是唯一索引冲突
				if !xorm.DuplicateError(err) {
					l.Errorf("InviteFriendToGroup InsertMany error: %v", err)
					return err
				}
			}
			notice := &noticemodel.Notice{
				ConvId: pb.HiddenConvIdGroup(apply.GroupId),
				Options: noticemodel.NoticeOption{
					StorageForClient: false,
					UpdateConvNotice: false,
				},
				ContentType: pb.NoticeContentType_NewGroupMember,
				Content: utils.AnyToBytes(pb.NoticeContent_NewGroupMember{
					GroupId:  apply.GroupId,
					MemberId: apply.UserId,
				}),
				UniqueId: "member",
				Title:    "",
				Ext:      nil,
			}
			err = notice.Insert(l.ctx, tx, l.svcCtx.Redis())
			if err != nil {
				l.Errorf("insert notice failed, err: %v", err)
				return err
			}
			return nil
		}
		return nil
	}, func(tx *gorm.DB) error {
		for _, manager := range groupManagers {
			notice := &noticemodel.Notice{
				ConvId: pb.HiddenConvIdGroupMember(),
				UserId: manager.UserId,
				Options: noticemodel.NoticeOption{
					StorageForClient: false,
					UpdateConvNotice: false,
				},
				ContentType: pb.NoticeContentType_ApplyToBeGroupMember,
				Content: utils.AnyToBytes(pb.NoticeContent_ApplyToBeGroupMember{
					ApplyId:      apply.Id,
					GroupId:      apply.GroupId,
					UserId:       apply.UserId,
					Result:       apply.Result,
					Reason:       apply.Reason,
					ApplyTime:    apply.ApplyTime,
					HandleTime:   apply.HandleTime,
					HandleUserId: apply.HandleUserId,
				}),
				UniqueId: apply.Id,
				Title:    "",
				Ext:      nil,
			}
			err := notice.Insert(l.ctx, tx, l.svcCtx.Redis())
			if err != nil {
				l.Errorf("insert notice failed, err: %v", err)
				return err
			}
			return nil
		}
		return nil
	})
	if err != nil {
		l.Errorf("HandleGroupApply Transaction error: %v", err)
		return &pb.HandleGroupApplyResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	utils.RetryProxy(context.Background(), 12, 1*time.Second, func() error {
		// 删除缓存
		{
			err = groupmodel.FlushGroupsByUserIdCache(l.ctx, l.svcCtx.Redis(), apply.UserId)
			if err != nil {
				l.Errorf("InviteFriendToGroup FlushGroupsByUserIdCache error: %v", err)
				return err
			}
		}
		// 预热缓存
		go xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "CacheWarmUp", func(ctx context.Context) {
			_, err = groupmodel.ListGroupsByUserIdFromMysql(ctx, l.svcCtx.Mysql(), l.svcCtx.Redis(), apply.UserId)
		}, propagation.MapCarrier{
			"group_id": apply.GroupId,
		})
		_, err := l.svcCtx.MsgService().FlushUsersSubConv(l.ctx, &pb.FlushUsersSubConvReq{UserIds: []string{apply.UserId}})
		if err != nil {
			l.Errorf("FlushUsersSubConv failed, err: %v", err)
			return err
		}
		if in.Result == pb.GroupApplyHandleResult_AGREE {
			_, err = l.svcCtx.NoticeService().GetUserNoticeData(l.ctx, &pb.GetUserNoticeDataReq{
				CommonReq: in.CommonReq,
				ConvId:    pb.HiddenConvIdGroup(apply.GroupId),
			})
			if err != nil {
				l.Errorf("SendNoticeData failed, err: %v", err)
			}
			_, err = NewSyncGroupMemberCountLogic(l.ctx, l.svcCtx).SyncGroupMemberCount(&pb.SyncGroupMemberCountReq{
				CommonReq: in.GetCommonReq(),
				GroupId:   apply.GroupId,
			})
			if err != nil {
				l.Errorf("SyncGroupMemberCount failed, err: %v", err)
				return err
			}
		}
		return nil
	})
	// 通知给群里所有的管理员
	go xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "SendNotice", func(ctx context.Context) {
		utils.RetryProxy(ctx, 12, time.Second, func() error {
			for _, manager := range groupManagers {
				_, err := l.svcCtx.NoticeService().GetUserNoticeData(ctx, &pb.GetUserNoticeDataReq{
					CommonReq: in.CommonReq,
					UserId:    manager.UserId,
					ConvId:    pb.HiddenConvIdGroupMember(),
				})
				if err != nil {
					l.Errorf("ApplyToBeGroupMember SendNoticeData error: %v", err)
					return err
				}
			}
			return nil
		})
	}, propagation.MapCarrier{
		"groupId":  apply.GroupId,
		"userId":   in.CommonReq.UserId,
		"noticeId": apply.Id,
	})
	// 申请人发一条消息 自我介绍
	go xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "SendMyInfo", func(ctx context.Context) {
		if in.Result == pb.GroupApplyHandleResult_AGREE {
			// 获取我自己的信息
			var sender *pb.UserBaseInfo
			{
				userBaseInfo, err := l.svcCtx.UserService().BatchGetUserBaseInfo(ctx, &pb.BatchGetUserBaseInfoReq{
					CommonReq: in.CommonReq,
					Ids:       []string{apply.UserId},
				})
				if err != nil {
					l.Errorf("SendMsg failed, err: %v", err)
					return
				}
				if len(userBaseInfo.UserBaseInfos) == 0 {
					l.Errorf("SendMsg failed, err: %v", err)
					return
				}
				sender = userBaseInfo.UserBaseInfos[0]
			}
			_, err := msgservice.SendMsgSync(l.svcCtx.MsgService(), ctx, []*pb.MsgData{
				msgmodel.CreateTextMsgToGroup(
					sender,
					apply.GroupId,
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
				).ToMsgData(),
			})
			if err != nil {
				l.Errorf("SendMsg failed, err: %v", err)
			}
		}
	}, propagation.MapCarrier{
		"groupId":  apply.GroupId,
		"userId":   in.CommonReq.UserId,
		"noticeId": apply.Id,
	})
	return &pb.HandleGroupApplyResp{}, nil
}
