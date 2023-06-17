package smsservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/third/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SmsCodeVerifyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSmsCodeVerifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SmsCodeVerifyLogic {
	return &SmsCodeVerifyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SmsCodeVerify 验证短信
func (l *SmsCodeVerifyLogic) SmsCodeVerify(in *pb.SmsCodeVerifyReq) (*pb.SmsCodeVerifyResp, error) {
	// todo: add your logic here and delete this line

	return &pb.SmsCodeVerifyResp{Success: true}, nil
}
