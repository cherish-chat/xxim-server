package thirdservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/third/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

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
func (l *CaptchaVerifyLogic) CaptchaVerify(in *pb.CaptchaVerifyReq) (*pb.CaptchaVerifyResp, error) {
	// todo: add your logic here and delete this line

	return &pb.CaptchaVerifyResp{}, nil
}
