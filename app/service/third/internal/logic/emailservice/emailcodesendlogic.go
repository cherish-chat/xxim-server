package emailservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/third/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type EmailCodeSendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEmailCodeSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmailCodeSendLogic {
	return &EmailCodeSendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// EmailCodeSend 发送邮件
func (l *EmailCodeSendLogic) EmailCodeSend(in *peerpb.EmailCodeSendReq) (*peerpb.EmailCodeSendResp, error) {
	// todo: add your logic here and delete this line

	return &peerpb.EmailCodeSendResp{}, nil
}
