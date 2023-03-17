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
	latestGetConvIdsTTl := 0
	if len(latestGetConvIds) > 0 {
		// 获取这个key的过期秒数
		ttl, err := l.svcCtx.Redis().TtlCtx(l.ctx, latestGetConvIdsRedisKey)
		if err != nil {
			l.Errorf("ttl error: %v", err)
			return err
		}
		if ttl <= 0 {
			// 说明这个key永不过期
			ttl = 60 * 30 // 30分钟
		}
		latestGetConvIdsTTl = ttl
	}
	// 对比两次convIds  获取应该删除的  和 应该添加的
	var delConvIds []string
	for _, id := range latestGetConvIds {
		var in bool
		for _, id2 := range convIds {
			if id == id2 {
				in = true
				break
			}
		}
		if !in {
			delConvIds = append(delConvIds, id)
		}
	}
	// 取消上次会话的订阅
	if len(delConvIds) > 0 {
		var keys []string
		for _, id := range delConvIds {
			keys = append(keys, rediskey.ConvMembersSubscribed(id))
		}
		err := xredis.MZRem(l.svcCtx.Redis(), l.ctx, keys, rediskey.ConvMemberPodIp(userId))
		if err != nil {
			l.Errorf("mzrem error: %v", err)
			return err
		}
	}

	// 更新最新的会话ids
	if len(convIds) > 0 {
		err := l.svcCtx.Redis().SetexCtx(l.ctx, latestGetConvIdsRedisKey, strings.Join(convIds, ","), latestGetConvIdsTTl)
		if err != nil {
			l.Errorf("setex error: %v", err)
			return err
		}
	}

	// 新增订阅
	if len(convIds) > 0 {
		var keys []string
		for _, id := range convIds {
			keys = append(keys, rediskey.ConvMembersSubscribed(id))
		}
		err := xredis.MZAddEx(l.svcCtx.Redis(), l.ctx, keys, time.Now().UnixMilli(), rediskey.ConvMemberPodIp(userId), 60*60*24*30)
		if err != nil {
			l.Errorf("mzadd error: %v", err)
			return err
		}
	}

	return nil
}
