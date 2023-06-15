package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ConsumerLogic struct {
	svcCtx *svc.ServiceContext
	ctx    context.Context
	logx.Logger
}

func NewConsumerLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConsumerLogic {
	return &ConsumerLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

// AfterRegister TODO: 补充注册后的逻辑
func (l *ConsumerLogic) AfterRegister(ctx context.Context, topic string, msg []byte) error {
	l.Infof("topic: %s, msg: %s", topic, string(msg))
	return nil
}
