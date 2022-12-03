package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/notice/noticemodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"go.opentelemetry.io/otel/propagation"
	"time"

	"github.com/cherish-chat/xxim-server/app/notice/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AckNoticeDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAckNoticeDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AckNoticeDataLogic {
	return &AckNoticeDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// AckNoticeData 确认通知数据
func (l *AckNoticeDataLogic) AckNoticeData(in *pb.AckNoticeDataReq) (*pb.AckNoticeDataResp, error) {
	if in.Success {
		notice := &noticemodel.Notice{}
		err := l.svcCtx.Mysql().Model(notice).Where("noticeId = ?", in.NoticeId).First(notice).Error
		if err != nil {
			l.Errorf("find notice error: %v", err)
			return &pb.AckNoticeDataResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
		err = l.svcCtx.Redis().HsetCtx(l.ctx, rediskey.UserAckRecord(in.CommonReq.UserId, in.CommonReq.DeviceId), notice.ConvId, utils.AnyToString(notice.CreateTime))
		if err != nil {
			l.Errorf("redis hset error: %v", err)
			return &pb.AckNoticeDataResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	// 再次查询
	go xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "GetUserNoticeConvIds", func(ctx context.Context) {
		for i := 0; i < 12; i++ {
			resp, err := NewGetUserNoticeDataLogic(ctx, l.svcCtx).GetUserNoticeData(&pb.GetUserNoticeDataReq{
				CommonReq: in.CommonReq,
				UserId:    in.CommonReq.UserId,
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
	return &pb.AckNoticeDataResp{}, nil
}
