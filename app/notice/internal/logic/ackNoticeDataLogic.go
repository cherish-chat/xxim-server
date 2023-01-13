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
func (l *AckNoticeDataLogic) ackNoticeData(in *pb.AckNoticeDataReq) (*pb.AckNoticeDataResp, error) {
	notice := &noticemodel.Notice{}
	err := l.svcCtx.Mysql().Model(&noticemodel.Notice{}).Where("noticeId = ? AND convId = ?", in.NoticeId, in.ConvId).Limit(1).Find(notice).Error
	if err != nil {
		l.Errorf("find notice error: %v", err)
		return &pb.AckNoticeDataResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	hmSetMap := map[string]string{
		in.ConvId: utils.AnyToString(notice.CreateTime),
	}
	err = l.svcCtx.Redis().HmsetCtx(l.ctx, rediskey.UserAckRecord(in.CommonReq.UserId, in.CommonReq.DeviceId), hmSetMap)
	if err != nil {
		l.Errorf("redis hset error: %v", err)
		return &pb.AckNoticeDataResp{CommonResp: pb.NewRetryErrorResp()}, err
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

// AckNoticeData 确认通知数据 TODO 临时使用
func (l *AckNoticeDataLogic) AckNoticeData(in *pb.AckNoticeDataReq) (*pb.AckNoticeDataResp, error) {
	if in.ConvId == "" {
		notices := make([]*noticemodel.Notice, 0)
		err := l.svcCtx.Mysql().Model(&noticemodel.Notice{}).
			Where("noticeId = ? AND userId = ?", in.NoticeId, in.CommonReq.UserId).Find(&notices).Error
		if err != nil {
			l.Errorf("find notice error: %v", err)
			return &pb.AckNoticeDataResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
		convNoticesMap := make(map[string][]*noticemodel.Notice)
		convMaxCreateTimeMap := make(map[string]int64)
		for _, notice := range notices {
			if _, ok := convNoticesMap[notice.ConvId]; !ok {
				convNoticesMap[notice.ConvId] = make([]*noticemodel.Notice, 0)
			}
			convNoticesMap[notice.ConvId] = append(convNoticesMap[notice.ConvId], notice)
			if _, ok := convMaxCreateTimeMap[notice.ConvId]; !ok {
				convMaxCreateTimeMap[notice.ConvId] = notice.CreateTime
			}
			if notice.CreateTime > convMaxCreateTimeMap[notice.ConvId] {
				convMaxCreateTimeMap[notice.ConvId] = notice.CreateTime
			}
		}
		hmSetMap := make(map[string]string)
		for convId, createTime := range convMaxCreateTimeMap {
			hmSetMap[convId] = utils.AnyToString(createTime)
		}
		err = l.svcCtx.Redis().HmsetCtx(l.ctx, rediskey.UserAckRecord(in.CommonReq.UserId, in.CommonReq.DeviceId), hmSetMap)
		if err != nil {
			l.Errorf("redis hset error: %v", err)
			return &pb.AckNoticeDataResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	} else {
		return l.ackNoticeData(in)
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
