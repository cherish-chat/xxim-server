package ws

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/im/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ConsumePushLogic struct {
	svcCtx *svc.ServiceContext
	ctx    context.Context
	logx.Logger
}

func NewConsumePushLogic(svcCtx *svc.ServiceContext, ctx context.Context) *ConsumePushLogic {
	return &ConsumePushLogic{svcCtx: svcCtx, ctx: ctx, Logger: logx.WithContext(ctx)}
}

func (l *ConsumePushLogic) ConsumePush(key string, payload []byte) error {
	// TODO add your logic here and delete this line
	return nil
}
