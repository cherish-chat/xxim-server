package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xredis"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"strings"
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
		err := l.SetUserSubscriptions(userId, in.CompareConvIds)
		if err != nil {
			l.Errorf("set user subscriptions error: %v", err)
			return pb.NewRetryErrorResp(), err
		}
	}
	return &pb.CommonResp{}, nil
}

func (l *FlushUsersSubConvLogic) SetUserSubscriptions(userId string, compareConvIds []string) error {
	var convIds []string
	convIdOfUser, err := l.svcCtx.ImService().GetAllConvIdOfUser(l.ctx, &pb.GetAllConvIdOfUserReq{
		UserId: userId,
	})
	if err != nil {
		l.Errorf("get all conv id of user error: %v", err)
		return err
	}
	convIds = convIdOfUser.ConvIds
	convIdMap := make(map[string]bool)
	for _, id := range convIds {
		convIdMap[id] = true
	}
	for _, id := range compareConvIds {
		// 是否存在
		if _, ok := convIdMap[id]; !ok {
			// 如果是单聊
			if pb.IsSingleConv(id) {
				// 也订阅
				convIds = append(convIds, id)
			}
		}
	}
	// 获取上次订阅的会话ids
	latestGetConvIdsRedisKey := rediskey.LatestGetConvIds(userId)
	val, _ := l.svcCtx.Redis().GetCtx(l.ctx, latestGetConvIdsRedisKey)
	latestGetConvIds := strings.Split(val, ",")
	// 取消上次会话的订阅
	if len(latestGetConvIds) > 0 {
		var keys []string
		for _, id := range latestGetConvIds {
			keys = append(keys, rediskey.ConvMembersSubscribed(id))
		}
		err := xredis.MZRem(l.svcCtx.RedisSub(), l.ctx, keys, rediskey.ConvMemberPodIp(userId))
		if err != nil {
			l.Errorf("mzrem error: %v", err)
			return err
		}
	}
	// mzadd and setex
	if len(convIds) > 0 {
		var keys []string
		for _, id := range convIds {
			keys = append(keys, rediskey.ConvMembersSubscribed(id))
		}
		err := xredis.MZAddEx(l.svcCtx.RedisSub(), l.ctx, keys, time.Now().UnixMilli(), rediskey.ConvMemberPodIp(userId), 60*5)
		if err != nil {
			l.Errorf("mzaddex error: %v", err)
			return err
		}
	}
	return nil
}
