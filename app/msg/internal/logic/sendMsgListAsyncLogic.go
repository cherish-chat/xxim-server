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

type SendMsgListAsyncLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendMsgListAsyncLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendMsgListAsyncLogic {
	return &SendMsgListAsyncLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendMsgListAsyncLogic) SendMsgListAsync(in *pb.SendMsgListReq) (*pb.SendMsgListResp, error) {
	// 发送到消息队列
	var options []xtdmq.ProducerOptFunc
	if in.DeliverAfter != nil {
		options = append(options, xtdmq.ProduceWithDeliverAfter(time.Second*time.Duration(*in.DeliverAfter)))
	}
	_, err := l.svcCtx.MsgProducer().Produce(l.ctx, "msg", utils.ProtoToBytes(&pb.MsgMQBody{
		Event: pb.MsgMQBody_SendMsgListSync,
		Data:  utils.ProtoToBytes(in),
	}), options...)
	if err != nil {
		l.Errorf("MsgProducer.Produce error: %v", err)
		return &pb.SendMsgListResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.SendMsgListResp{}, nil
}
