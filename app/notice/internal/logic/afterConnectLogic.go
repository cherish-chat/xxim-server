package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xredis"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"go.opentelemetry.io/otel/propagation"
	"time"

	"github.com/cherish-chat/xxim-server/app/notice/internal/svc"
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

// AfterConnect conn hook
func (l *AfterConnectLogic) AfterConnect(in *pb.AfterConnectReq) (*pb.CommonResp, error) {
	err := l.SetUserSubscriptions(in.ConnParam.UserId)
	if err != nil {
		return &pb.CommonResp{}, err
	}
	go xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "GetUserNoticeConvIds", func(ctx context.Context) {
		for i := 0; i < 12; i++ {
			resp, err := NewGetUserNoticeDataLogic(ctx, l.svcCtx).GetUserNoticeData(&pb.GetUserNoticeDataReq{
				UserId: in.ConnParam.UserId,
			})
			if err != nil {
				l.Errorf("get user notice data error: %v", err)
			} else if resp.CommonResp != nil && resp.CommonResp.Failed() {
				l.Errorf("get user notice data failed: %v", utils.AnyToString(resp))
			} else {
				break
			}
			time.Sleep(2 * time.Second)
		}
	}, propagation.MapCarrier{})
	return &pb.CommonResp{}, nil
}

func (l *AfterConnectLogic) SetUserSubscriptions(userId string) error {
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
