package ws

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/im/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ConsumeInternalLogic struct {
	svcCtx *svc.ServiceContext
	ctx    context.Context
	logx.Logger
}

func NewConsumeInternalLogic(svcCtx *svc.ServiceContext, ctx context.Context) *ConsumeInternalLogic {
	return &ConsumeInternalLogic{svcCtx: svcCtx, ctx: ctx, Logger: logx.WithContext(ctx)}
}

func (l *ConsumeInternalLogic) ConsumeInternal(key string, payload []byte) error {
	// TODO add your logic here and delete this line
	return nil
}
