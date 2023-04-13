package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/utils/captcha"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginMSCaptchaLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginMSCaptchaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginMSCaptchaLogic {
	return &LoginMSCaptchaLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginMSCaptchaLogic) LoginMSCaptcha(in *pb.LoginMSCaptchaReq) (*pb.LoginMSCaptchaResp, error) {
	// 随机生成6位验证码
	code := utils.Random.String([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}, 6)
	var minute = 5
	// 保存验证码到redis
	captchaId := utils.GenId()
	key := rediskey.CaptchaCodeKey("LoginMS", captchaId)
	err := l.svcCtx.Redis().SetexCtx(l.ctx, key, code, minute*60)
	if err != nil {
		l.Errorf("SendSms failed: %v", err)
		return &pb.LoginMSCaptchaResp{CommonResp: pb.NewInternalErrorResp(err.Error())}, err
	}
	// 生成图片
	bytes := captcha.ImgText(300, 100, code)
	// base64
	b64 := utils.Base64.EncodeToString(bytes)
	return &pb.LoginMSCaptchaResp{CaptchaB64: b64, CaptchaId: captchaId}, nil
}
