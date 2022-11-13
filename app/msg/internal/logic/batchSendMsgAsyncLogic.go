package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xtdmq"
	"time"

	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchSendMsgAsyncLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchSendMsgAsyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchSendMsgAsyncLogic {
	return &BatchSendMsgAsyncLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BatchSendMsgAsyncLogic) BatchSendMsgAsync(in *pb.BatchSendMsgReq) (*pb.BatchSendMsgResp, error) {
	// 发送到消息队列
	var options []xtdmq.ProducerOptFunc
	if in.DeliverAfter != nil {
		options = append(options, xtdmq.ProduceWithDeliverAfter(time.Second*time.Duration(*in.DeliverAfter)))
	}
	_, err := l.svcCtx.MsgProducer().Produce(l.ctx, "msg", utils.ProtoToBytes(&pb.MsgMQBody{
		Event: pb.MsgMQBody_BatchSendMsgSync,
		Data:  utils.ProtoToBytes(in),
	}), options...)
	if err != nil {
		l.Errorf("MsgProducer.Produce error: %v", err)
		return &pb.BatchSendMsgResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.BatchSendMsgResp{}, nil
}
