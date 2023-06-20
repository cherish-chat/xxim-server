package messageservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/message/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type MessagePushLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMessagePushLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MessagePushLogic {
	return &MessagePushLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// MessagePush 推送消息
func (l *MessagePushLogic) MessagePush(in *pb.MessagePushReq) (*pb.MessagePushResp, error) {
	// todo: add your logic here and delete this line
	l.Debugf("MessagePush 推送消息: %+v", in)
	return &pb.MessagePushResp{}, nil
}
