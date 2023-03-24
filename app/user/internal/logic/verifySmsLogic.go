package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type VerifySmsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewVerifySmsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *VerifySmsLogic {
	return &VerifySmsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *VerifySmsLogic) VerifySms(in *pb.VerifySmsReq) (*pb.VerifySmsResp, error) {
	mobile := in.Phone
	if in.CountryCode != nil && *in.CountryCode != "" {
		mobile = *in.CountryCode + mobile
	} else {
		mobile = "+86" + mobile
	}
	smsCodeErrorKey := rediskey.SmsCodeErrorKey(in.Scene, mobile)
	// 验证码错误次数
	smsCodeErrorCount, _ := l.svcCtx.Redis().GetCtx(l.ctx, smsCodeErrorKey)
	if utils.AnyToInt64(smsCodeErrorCount) >= l.svcCtx.ConfigMgr.SmsErrorMaxCount(l.ctx) {
		return &pb.VerifySmsResp{CommonResp: pb.NewToastErrorResp("验证码错误次数过多，请稍后再试")}, nil
	}
	// 从redis中获取验证码
	key := rediskey.SmsCodeKey(in.Scene, mobile)
	code, err := l.svcCtx.Redis().GetCtx(l.ctx, key)
	if err != nil {
		l.Errorf("VerifySms failed: %v", err)
		return &pb.VerifySmsResp{CommonResp: pb.NewInternalErrorResp(err.Error())}, err
	}
	// 是否立刻删除
	if in.Delete {
		_, err = l.svcCtx.Redis().DelCtx(l.ctx, key)
		if err != nil {
			l.Errorf("VerifySms failed: %v", err)
			return &pb.VerifySmsResp{CommonResp: pb.NewInternalErrorResp(err.Error())}, err
		}
	}
	if code != in.Code {
		// 记录错误次数
		_, _ = l.svcCtx.Redis().IncrCtx(l.ctx, smsCodeErrorKey)
		// 设置过期时间
		_ = l.svcCtx.Redis().ExpireCtx(l.ctx, smsCodeErrorKey, int(l.svcCtx.ConfigMgr.SmsErrorPeriod(l.ctx)))
		return &pb.VerifySmsResp{CommonResp: pb.NewToastErrorResp("验证码错误")}, nil
	}
	return &pb.VerifySmsResp{}, nil
}
