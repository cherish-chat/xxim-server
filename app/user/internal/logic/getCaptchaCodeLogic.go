package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/utils/captcha"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCaptchaCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCaptchaCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCaptchaCodeLogic {
	return &GetCaptchaCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCaptchaCodeLogic) GetCaptchaCode(in *pb.GetCaptchaCodeReq) (*pb.GetCaptchaCodeResp, error) {
	// 随机生成6位验证码
	code := utils.Random.String([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}, 6)
	var minute = 5
	if in.ExpireMinute != nil && *in.ExpireMinute > 0 {
		minute = int(*in.ExpireMinute)
	}
	// 保存验证码到redis
	key := rediskey.CaptchaCodeKey(in.Scene, in.CommonReq.DeviceId)
	err := l.svcCtx.Redis().SetexCtx(l.ctx, key, code, minute*60)
	if err != nil {
		l.Errorf("SendSms failed: %v", err)
		return &pb.GetCaptchaCodeResp{CommonResp: pb.NewInternalErrorResp(err.Error())}, err
	}
	// 生成图片
	bytes := captcha.ImgText(300, 100, code)
	return &pb.GetCaptchaCodeResp{Captcha: bytes}, nil
}
