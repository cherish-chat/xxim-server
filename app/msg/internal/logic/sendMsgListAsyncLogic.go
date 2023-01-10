package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xtdmq"
	"github.com/cherish-chat/xxim-server/common/xtrace"
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
	// 判断每秒的请求量，如果达到阈值，走消息队列，否则直接调用SendMsgListSync
	allow := l.svcCtx.SyncSendMsgLimiter.AllowCtx(l.ctx)
	if !allow {
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
	} else {
		var resp *pb.SendMsgListResp
		var err error
		xtrace.StartFuncSpan(l.ctx, "SendMsgListSync", func(ctx context.Context) {
			resp, err = NewSendMsgListSyncLogic(ctx, l.svcCtx).SendMsgListSync(in)
		})
		return resp, err
	}
}
