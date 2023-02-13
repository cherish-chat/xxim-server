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
	return &pb.SendSmsResp{}, nil
}
