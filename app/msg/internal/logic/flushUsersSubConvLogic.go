package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xredis"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"time"

	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FlushUsersSubConvLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFlushUsersSubConvLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FlushUsersSubConvLogic {
	return &FlushUsersSubConvLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// FlushUsersSubConv 刷新用户订阅的会话
func (l *FlushUsersSubConvLogic) FlushUsersSubConv(in *pb.FlushUsersSubConvReq) (*pb.CommonResp, error) {
	for _, userId := range utils.Set(in.UserIds) {
		err := l.SetUserSubscriptions(userId)
		if err != nil {
			l.Errorf("set user subscriptions error: %v", err)
			return pb.NewRetryErrorResp(), err
		}
	}
	return &pb.CommonResp{}, nil
}

func (l *FlushUsersSubConvLogic) SetUserSubscriptions(userId string) error {
	var friendIds []string
	var groupIds []string
	var convIds []string
	// 获取用户订阅的好友列表
	{
		getFriendList, err := l.svcCtx.RelationService().GetFriendList(l.ctx, &pb.GetFriendListReq{
			CommonReq: &pb.CommonReq{
				UserId: userId,
			},
			Page: &pb.Page{
				Page: 1,
				Size: 0,
			},
			Opt: pb.GetFriendListReq_OnlyId,
		})
		if err != nil {
			l.Errorf("get friend list error: %v", err)
			return err
		}
		friendIds = getFriendList.Ids
		for _, id := range friendIds {
			convIds = append(convIds, pb.SingleConvId(userId, id))
		}
	}
	// 获取用户订阅的群组列表
	{
		getMyGroupList, err := l.svcCtx.GroupService().GetMyGroupList(l.ctx, &pb.GetMyGroupListReq{
			CommonReq: &pb.CommonReq{
				UserId: userId,
			},
			Page: &pb.Page{Page: 1},
			Filter: &pb.GetMyGroupListReq_Filter{
				FilterFold:   true,
				FilterShield: true,
			},
			Opt: pb.GetMyGroupListReq_ONLY_ID,
		})
		if err != nil {
			l.Errorf("get group list error: %v", err)
			return err
		}
		groupIds = getMyGroupList.Ids
		for _, id := range groupIds {
			convIds = append(convIds, id)
		}
	}
	// mzadd and setex
	if len(convIds) > 0 {
		var keys []string
		for _, id := range convIds {
			keys = append(keys, rediskey.ConvMembersSubscribed(id))
		}
		err := xredis.MZAddEx(l.svcCtx.Redis(), l.ctx, keys, time.Now().UnixMilli(), rediskey.ConvMemberPodIp(userId), 60*60*24)
		if err != nil {
			l.Errorf("mzaddex error: %v", err)
			return err
		}
	}
	return nil
}
