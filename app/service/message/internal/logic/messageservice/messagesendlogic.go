package messageservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xmq"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/message/internal/svc"

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
func (l *MessageSendLogic) MessageSend(in *peerpb.MessageSendReq) (*peerpb.MessageSendResp, error) {
	var (
		resp *peerpb.MessageSendResp
		err  error
	)
	{
		resp, err = l.validate(in)
		if err != nil {
			return resp, err
		}
		if !resp.Success {
			return resp, nil
		}
	}
	if !in.DisableQueue {
		err := l.svcCtx.MQ.Produce(context.Background(), xmq.TopicMessageInsert, utils.Proto.Marshal(&peerpb.MessageInsertReq{
			Header:   in.Header,
			Messages: []*peerpb.Message{in.Message},
		}))
		if err != nil {
			l.Errorf("produce message error: %v", err)
			return nil, err
		}
		return &peerpb.MessageSendResp{}, nil
	}

	_, _ = NewMessageInsertLogic(context.Background(), l.svcCtx).MessageInsert(&peerpb.MessageInsertReq{
		Header:   in.Header,
		Messages: []*peerpb.Message{in.Message},
	})
	return &peerpb.MessageSendResp{}, nil
}

func (l *MessageSendLogic) validate(in *peerpb.MessageSendReq) (*peerpb.MessageSendResp, error) {
	// TODO: add your logic here and delete this line
	return &peerpb.MessageSendResp{
		Success: true,
	}, nil
}
