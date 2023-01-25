package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/notice/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type TimerCleanSubscriptionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewTimerCleanSubscriptionLogic(svcCtx *svc.ServiceContext) *TimerCleanSubscriptionLogic {
	l := &TimerCleanSubscriptionLogic{svcCtx: svcCtx, ctx: context.Background()}
	l.Logger = logx.WithContext(l.ctx)
	return l
}

func (l *TimerCleanSubscriptionLogic) Start() {

}
