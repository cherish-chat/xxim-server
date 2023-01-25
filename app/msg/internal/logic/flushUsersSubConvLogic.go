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
	var convIds []string
	convIdOfUser, err := l.svcCtx.ImService().GetAllConvIdOfUser(l.ctx, &pb.GetAllConvIdOfUserReq{
		UserId: userId,
	})
	if err != nil {
		l.Errorf("get all conv id of user error: %v", err)
		return err
	}
	convIds = convIdOfUser.ConvIds
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
