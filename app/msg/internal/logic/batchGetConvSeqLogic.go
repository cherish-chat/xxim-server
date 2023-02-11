package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"go.opentelemetry.io/otel/propagation"

	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchGetConvSeqLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchGetConvSeqLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchGetConvSeqLogic {
	return &BatchGetConvSeqLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BatchGetConvSeq 批量获取会话的seq
func (l *BatchGetConvSeqLogic) BatchGetConvSeq(in *pb.BatchGetConvSeqReq) (*pb.BatchGetConvSeqResp, error) {
	go xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "KeepAlive", func(ctx context.Context) {
		_, e := l.svcCtx.ImService().KeepAlive(ctx, &pb.KeepAliveReq{CommonReq: in.GetCommonReq()})
		if e != nil {
			l.Errorf("KeepAlive failed, err: %v", e)
		}
	}, propagation.MapCarrier{
		"userId": in.CommonReq.UserId,
	})
	if len(in.ConvIdList) == 0 {
		return &pb.BatchGetConvSeqResp{}, nil
	}
	convMaxSeq, err := BatchGetConvMaxSeq(l.svcCtx.Redis(), l.ctx, in.CommonReq.UserId, in.ConvIdList)
	if err != nil {
		l.Errorf("BatchGetConvSeq BatchGetConvMaxSeq error: %v", err)
		return &pb.BatchGetConvSeqResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var resp = make(map[string]*pb.BatchGetConvSeqResp_ConvSeq)
	for convId, seq := range convMaxSeq {
		resp[convId] = &pb.BatchGetConvSeqResp_ConvSeq{
			ConvId:     convId,
			MinSeq:     utils.AnyToString(seq.minSeq),
			MaxSeq:     utils.AnyToString(seq.maxSeq),
			UpdateTime: utils.AnyToString(seq.updateTime),
		}
	}
	for _, convId := range in.ConvIdList {
		if _, ok := resp[convId]; !ok {
			resp[convId] = &pb.BatchGetConvSeqResp_ConvSeq{
				ConvId:     convId,
				MinSeq:     "0",
				MaxSeq:     "0",
				UpdateTime: "0",
			}
		}
	}
	return &pb.BatchGetConvSeqResp{ConvSeqMap: resp}, nil
}
