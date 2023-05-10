package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifyLoginMSCaptchaCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifyLoginMSCaptchaCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifyLoginMSCaptchaCodeLogic {
	return &VerifyLoginMSCaptchaCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *VerifyLoginMSCaptchaCodeLogic) VerifyLoginMSCaptchaCode(in *pb.VerifyLoginMSCaptchaCodeReq) (*pb.VerifyLoginMSCaptchaCodeResp, error) {
	key := rediskey.CaptchaCodeKey("LoginMS", in.CaptchaId)
	code, err := l.svcCtx.Redis().GetCtx(l.ctx, key)
	if err != nil {
		l.Errorf("VerifyCaptchaCode failed: %v", err)
		return &pb.VerifyLoginMSCaptchaCodeResp{CommonResp: pb.NewInternalErrorResp(err.Error())}, err
	}
	if code == "" {
		l.Errorf("VerifyCaptchaCode failed: %v", err)
		return &pb.VerifyLoginMSCaptchaCodeResp{CommonResp: pb.NewInternalErrorResp(l.svcCtx.T(in.CommonReq.Language, "图形验证码已失效"))}, err
	}
	// 是否立刻删除
	if in.Delete {
		_, err = l.svcCtx.Redis().DelCtx(l.ctx, key)
		if err != nil {
			l.Errorf("VerifySms failed: %v", err)
			return &pb.VerifyLoginMSCaptchaCodeResp{CommonResp: pb.NewInternalErrorResp(err.Error())}, err
		}
	}
	if code != in.Code {
		l.Infof("code verify failed: %s != %s", code, in.Code)
		return &pb.VerifyLoginMSCaptchaCodeResp{CommonResp: pb.NewToastErrorResp(l.svcCtx.T(in.CommonReq.Language, "图形验证码错误"))}, nil
	}
	return &pb.VerifyLoginMSCaptchaCodeResp{CommonResp: pb.NewSuccessResp()}, nil
}
