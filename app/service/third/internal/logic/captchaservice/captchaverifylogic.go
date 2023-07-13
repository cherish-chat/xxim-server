package captchaservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/third/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CaptchaVerifyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCaptchaVerifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CaptchaVerifyLogic {
	return &CaptchaVerifyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CaptchaVerify 验证图形验证码
func (l *CaptchaVerifyLogic) CaptchaVerify(in *peerpb.CaptchaVerifyReq) (*peerpb.CaptchaVerifyResp, error) {
	// todo: add your logic here and delete this line

	return &peerpb.CaptchaVerifyResp{}, nil
}
