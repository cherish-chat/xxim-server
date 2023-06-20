package messageservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/message/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MessageSendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMessageSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessageSendLogic {
	return &MessageSendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// MessageSend 发送消息
func (l *MessageSendLogic) MessageSend(in *pb.MessageSendReq) (*pb.MessageSendResp, error) {
	messageBatchSendResp, err := NewMessageBatchSendLogic(l.ctx, l.svcCtx).MessageBatchSend(&pb.MessageBatchSendReq{
		Header:       in.Header,
		Messages:     []*pb.Message{in.Message},
		DisableQueue: in.DisableQueue,
	})
	var header *pb.ResponseHeader
	if messageBatchSendResp != nil {
		header = messageBatchSendResp.GetHeader()
	}
	return &pb.MessageSendResp{
		Header: header,
	}, err
}
