package emailservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/third/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type EmailCodeVerifyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewEmailCodeVerifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EmailCodeVerifyLogic {
	return &EmailCodeVerifyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// EmailCodeVerify 验证邮件
func (l *EmailCodeVerifyLogic) EmailCodeVerify(in *peerpb.EmailCodeVerifyReq) (*peerpb.EmailCodeVerifyResp, error) {
	// todo: add your logic here and delete this line

	return &peerpb.EmailCodeVerifyResp{Success: true}, nil
}
