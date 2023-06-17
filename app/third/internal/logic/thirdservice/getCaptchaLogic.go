package thirdservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/third/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCaptchaLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCaptchaLogic {
	return &GetCaptchaLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetCaptcha 获取图形验证码
func (l *GetCaptchaLogic) GetCaptcha(in *pb.GetCaptchaReq) (*pb.GetCaptchaResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetCaptchaResp{}, nil
}
