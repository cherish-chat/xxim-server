package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"strconv"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendSmsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendSmsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendSmsLogic {
	return &SendSmsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendSmsLogic) SendSms(in *pb.SendSmsReq) (*pb.SendSmsResp, error) {
	sender, err := l.svcCtx.SmsSender()
	if err != nil {
		l.Errorf("SendSms failed: %v", err)
		return &pb.SendSmsResp{CommonResp: pb.NewInternalErrorResp(err.Error())}, err
	}
	mobile := in.Phone
	if in.CountryCode != nil && *in.CountryCode != "" {
		mobile = *in.CountryCode + mobile
	}
	// 手机号限流是否触发
	if exist, _ := l.svcCtx.Redis().Exists(rediskey.SmsSendLimitKey(in.Scene, mobile)); exist {
		return &pb.SendSmsResp{CommonResp: pb.NewToastErrorResp("发送频率过快")}, nil
	}
	// 每天上限
	if val, _ := l.svcCtx.Redis().GetCtx(l.ctx, rediskey.SmsSendLimitEverydayKey(in.Scene, mobile)); val != "" {
		count := utils.AnyToInt64(val)
		if count >= l.svcCtx.ConfigMgr.SmsSendLimitEveryday(l.ctx) {
			return &pb.SendSmsResp{CommonResp: pb.NewToastErrorResp("发送次数超过上限")}, nil
		}
	}
	// 随机生成6位验证码
	code := utils.Random.String([]string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0"}, 6)
	var minute = 5
	if in.ExpireMinute != nil && *in.ExpireMinute > 0 {
		minute = int(*in.ExpireMinute)
	}
	// 保存验证码到redis
	key := rediskey.SmsCodeKey(in.Scene, mobile)
	err = l.svcCtx.Redis().SetexCtx(l.ctx, key, code, minute*60)
	if err != nil {
		l.Errorf("SendSms failed: %v", err)
		return &pb.SendSmsResp{CommonResp: pb.NewInternalErrorResp(err.Error())}, err
	}
	err = sender.SendMsg([]string{mobile}, code, strconv.Itoa(minute))
	if err != nil {
		l.Errorf("SendSms failed: %v", err)
		return &pb.SendSmsResp{CommonResp: pb.NewInternalErrorResp(err.Error())}, err
	}
	// 限流 n 秒内不能再次发送
	_ = l.svcCtx.Redis().Setex(rediskey.SmsSendLimitKey(in.Scene, mobile), "1", l.svcCtx.ConfigMgr.SmsSendLimitInterval(l.ctx))
	// SmsSendLimitEveryday
	_, _ = l.svcCtx.Redis().IncrbyCtx(l.ctx, rediskey.SmsSendLimitEverydayKey(in.Scene, mobile), 1)
	_ = l.svcCtx.Redis().ExpireCtx(l.ctx, rediskey.SmsSendLimitEverydayKey(in.Scene, mobile), 24*60*60)
	return &pb.SendSmsResp{}, nil
}
