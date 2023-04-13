package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyCaptchaCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifyCaptchaCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyCaptchaCodeLogic {
	return &VerifyCaptchaCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *VerifyCaptchaCodeLogic) VerifyCaptchaCode(in *pb.VerifyCaptchaCodeReq) (*pb.VerifyCaptchaCodeResp, error) {
	key := rediskey.CaptchaCodeKey(in.Scene, in.DeviceId)
	code, err := l.svcCtx.Redis().GetCtx(l.ctx, key)
	if err != nil {
		l.Errorf("VerifyCaptchaCode failed: %v", err)
		return &pb.VerifyCaptchaCodeResp{CommonResp: pb.NewInternalErrorResp(err.Error())}, err
	}
	// 是否立刻删除
	if in.Delete {
		_, err = l.svcCtx.Redis().DelCtx(l.ctx, key)
		if err != nil {
			l.Errorf("VerifySms failed: %v", err)
			return &pb.VerifyCaptchaCodeResp{CommonResp: pb.NewInternalErrorResp(err.Error())}, err
		}
	}
	if code != in.Code {
		return &pb.VerifyCaptchaCodeResp{CommonResp: pb.NewToastErrorResp("图形验证码错误")}, nil
	}
	return &pb.VerifyCaptchaCodeResp{}, nil
}
