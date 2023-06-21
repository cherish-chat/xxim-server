package messageservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xmq"

	"github.com/cherish-chat/xxim-server/app/message/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MessageBatchSendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMessageBatchSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessageBatchSendLogic {
	return &MessageBatchSendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// MessageBatchSend 批量发送消息
func (l *MessageBatchSendLogic) MessageBatchSend(in *pb.MessageBatchSendReq) (*pb.MessageBatchSendResp, error) {
	//验证
	var (
		resp *pb.MessageBatchSendResp
		err  error
	)
	{
		resp, err = l.validate(in)
		if err != nil {
			return resp, err
		}
		if resp.GetHeader().GetCode() != pb.ResponseCode_SUCCESS {
			return resp, nil
		}
	}
	//判断是否禁走队列
	if !in.DisableQueue {
		if !l.svcCtx.SendMsgTokenLimiter.AllowCtx(l.ctx) {
			err := l.svcCtx.MQ.Produce(l.ctx, xmq.TopicMessageBatchSend, utils.Json.MarshalToBytes(&pb.MessageInsertReq{
				Header:   in.Header,
				Messages: in.Messages,
			}))
			if err != nil {
				l.Errorf("produce message error: %v", err)
				return &pb.MessageBatchSendResp{}, err
			}
			return &pb.MessageBatchSendResp{}, nil
		}
	}

	//直接插入
	_, _ = NewMessageInsertLogic(l.ctx, l.svcCtx).MessageInsert(&pb.MessageInsertReq{
		Header:   in.Header,
		Messages: in.Messages,
	})
	return &pb.MessageBatchSendResp{}, nil
}

// validate 验证是否允许发送消息
func (l *MessageBatchSendLogic) validate(in *pb.MessageBatchSendReq) (*pb.MessageBatchSendResp, error) {
	// TODO: 验证是否允许发送消息
	return &pb.MessageBatchSendResp{}, nil
}
