package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/xredis"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"time"

	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AfterConnectLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAfterConnectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AfterConnectLogic {
	return &AfterConnectLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// conn hook
func (l *AfterConnectLogic) AfterConnect(in *pb.AfterConnectReq) (*pb.CommonResp, error) {
	err := l.SetUserSubscriptions(in.ConnParam.UserId, in.ConnParam.PodIp)
	if err != nil {
		return &pb.CommonResp{}, err
	}
	return &pb.CommonResp{}, nil
}

func (l *AfterConnectLogic) SetUserSubscriptions(userId string, podIp string) error {
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
		err := xredis.MZAddEx(l.svcCtx.Redis(), l.ctx, keys, time.Now().UnixMilli(), rediskey.ConvMemberPodIp(userId, podIp), 60*60*24)
		if err != nil {
			l.Errorf("mzaddex error: %v", err)
			return err
		}
	}
	return nil
}
