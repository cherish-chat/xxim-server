package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/utils"
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
	err := NewSetUserSubscriptionsLogic(l.ctx, l.svcCtx).SetUserSub(in.ConnParam.UserId)
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
