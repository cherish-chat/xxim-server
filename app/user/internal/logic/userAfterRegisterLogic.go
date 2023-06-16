package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type UserAfterRegisterLogic struct {
	svcCtx *svc.ServiceContext
	ctx    context.Context
	logx.Logger
}

func NewUserAfterRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAfterRegisterLogic {
	return &UserAfterRegisterLogic{ctx: ctx, svcCtx: svcCtx}
}

// AfterRegister TODO: 补充注册后的逻辑
func (l *UserAfterRegisterLogic) AfterRegister(topic string, msg []byte) error {
	l.Infof("topic: %s, msg: %s", topic, string(msg))
	return nil
}
