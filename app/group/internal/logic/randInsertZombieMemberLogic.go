package logic

import (
	"context"
	"fmt"
	"github.com/cherish-chat/xxim-server/app/group/groupmodel"
	msgservice "github.com/cherish-chat/xxim-server/app/msg/msgService"
	"github.com/cherish-chat/xxim-server/app/msg/msgmodel"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
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

type RandInsertZombieMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRandInsertZombieMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RandInsertZombieMemberLogic {
	return &RandInsertZombieMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RandInsertZombieMember 随机插入僵尸用户
func (l *RandInsertZombieMemberLogic) RandInsertZombieMember(in *pb.RandInsertZombieMemberReq) (*pb.RandInsertZombieMemberResp, error) {
	if in.Count == 0 {
		return &pb.RandInsertZombieMemberResp{}, nil
	}
	if in.Count > 10000 {
		in.Count = 10000
	}
	groupMemberModel := groupmodel.GroupMember{}
	// mysql rand get user
	var users []*usermodel.User
	var memberIds []string
	err := l.svcCtx.Mysql().Model(&usermodel.User{}).Where(
		fmt.Sprintf("role = ? AND id NOT IN (SELECT userId FROM %s WHERE groupId = ?)", groupMemberModel.TableName()),
		usermodel.RoleZombie, in.GroupId,
	).Order("RAND()").Limit(int(in.Count)).Find(&users).Error
	if err != nil {
		l.Errorf("RandInsertZombieMember error: %v", err)
		return &pb.RandInsertZombieMemberResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	if len(users) == 0 {
		return &pb.RandInsertZombieMemberResp{}, nil
	}
	err = xorm.Transaction(l.svcCtx.Mysql(), func(tx *gorm.DB) error {
		var members []*groupmodel.GroupMember
		for _, user := range users {
			members = append(members, &groupmodel.GroupMember{
				GroupId:    in.GroupId,
				UserId:     user.Id,
				CreateTime: time.Now().UnixMilli(),
				Role:       groupmodel.RoleType_MEMBER,
				Remark:     "",
				UnbanTime:  0,
			})
			memberIds = append(memberIds, user.Id)
		}
		err := tx.CreateInBatches(members, 100).Error
		if err != nil {
			l.Errorf("RandInsertZombieMember error: %v", err)
			return err
		}
		return nil
	}, func(tx *gorm.DB) error {
		return groupmodel.FlushGroupMemberListCache(l.ctx, l.svcCtx.Redis(), in.GroupId)
	})
	if err != nil {
		l.Errorf("RandInsertZombieMember error: %v", err)
		return &pb.RandInsertZombieMemberResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	utils.RetryProxy(context.Background(), 12, 1*time.Second, func() error {
		// 删除缓存
		{
			err = groupmodel.FlushGroupsByUserIdCache(l.ctx, l.svcCtx.Redis(), memberIds...)
			if err != nil {
				l.Errorf("InviteFriendToGroup FlushGroupsByUserIdCache error: %v", err)
				return err
			}
		}
		_, err := l.svcCtx.MsgService().FlushUsersSubConv(l.ctx, &pb.FlushUsersSubConvReq{UserIds: memberIds})
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
				Ids:       memberIds,
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
		for _, id := range memberIds {
			sender, ok := userBaseInfoMap[id]
			if !ok {
				continue
			}
			_, err := msgservice.SendMsgSync(l.svcCtx.MsgService(), ctx, []*pb.MsgData{
				msgmodel.CreateTextMsgToGroup(
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
				).ToMsgData(),
			})
			if err != nil {
				l.Errorf("SendMsg failed, err: %v", err)
			}
		}
	}, propagation.MapCarrier{
		"groupId": in.GroupId,
	})
	return &pb.RandInsertZombieMemberResp{}, nil
}
