package smsservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/third/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SmsCodeSendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSmsCodeSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SmsCodeSendLogic {
	return &SmsCodeSendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SmsCodeSend 发送短信
func (l *SmsCodeSendLogic) SmsCodeSend(in *peerpb.SmsCodeSendReq) (*peerpb.SmsCodeSendResp, error) {
	// todo: add your logic here and delete this line

	return &peerpb.SmsCodeSendResp{}, nil
}
