package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/xredis"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"time"

	"github.com/cherish-chat/xxim-server/app/notice/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetUserSubscriptionsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetUserSubscriptionsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetUserSubscriptionsLogic {
	return &SetUserSubscriptionsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SetUserSubscriptions 设置用户订阅
func (l *SetUserSubscriptionsLogic) SetUserSubscriptions(in *pb.SetUserSubscriptionsReq) (*pb.CommonResp, error) {
	for _, id := range in.UserIds {
		err := l.SetUserSub(id)
		if err != nil {
			return pb.NewRetryErrorResp(), nil
		}
	}
	return pb.NewSuccessResp(), nil
}

func (l *SetUserSubscriptionsLogic) SetUserSub(userId string) error {
	var convIds []string
	// 获取用户订阅的通知号列表
	{
		var getUserNoticeConvIdsResp *pb.GetUserNoticeConvIdsResp
		var err error
		xtrace.StartFuncSpan(l.ctx, "GetUserNoticeConvIds", func(ctx context.Context) {
			getUserNoticeConvIdsResp, err = NewGetUserNoticeConvIdsLogic(ctx, l.svcCtx).GetUserNoticeConvIds(&pb.GetUserNoticeConvIdsReq{
				UserId: userId,
			})
		})
		if err != nil {
			l.Errorf("get group list error: %v", err)
			return err
		}
		convIds = append(convIds, getUserNoticeConvIdsResp.ConvIds...)
	}
	// mzadd and setex
	if len(convIds) > 0 {
		var keys []string
		for _, id := range convIds {
			keys = append(keys, rediskey.NoticeConvMembersSubscribed(id))
		}
		err := xredis.MZAddEx(l.svcCtx.Redis(), l.ctx, keys, time.Now().UnixMilli(), rediskey.ConvMemberPodIp(userId), 60*60*24)
		if err != nil {
			l.Errorf("mzaddex error: %v", err)
			return err
		}
	}
	return nil
}
