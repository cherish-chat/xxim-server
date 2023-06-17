package svc

import (
	"github.com/cherish-chat/xxim-server/app/third/client/captchaservice"
	"github.com/cherish-chat/xxim-server/app/third/client/emailservice"
	"github.com/cherish-chat/xxim-server/app/third/client/smsservice"
	"github.com/cherish-chat/xxim-server/app/user/internal/config"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xcache"
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"github.com/cherish-chat/xxim-server/common/xmq"
	"github.com/qiniu/qmgo"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"time"
)

type ServiceContext struct {
	Config config.Config
	Redis  *redis.Redis

	UserCollection        *qmgo.QmgoClient
	UserSettingCollection *qmgo.QmgoClient

	MQ  xmq.MQ
	Jwt *utils.Jwt

	SmsService     smsservice.SmsService
	EmailService   emailservice.EmailService
	CaptchaService captchaservice.CaptchaService
}

func NewServiceContext(c config.Config) *ServiceContext {
	s := &ServiceContext{
		Config:                c,
		Redis:                 xcache.MustNewRedis(c.RedisConf),
		UserCollection:        xmgo.MustNewMongoCollection(c.MongoCollection.User, &usermodel.User{}),
		UserSettingCollection: xmgo.MustNewMongoCollection(c.MongoCollection.UserSetting, &usermodel.UserSetting{}),
	}

	//third rpc
	{
		thirdClient := zrpc.MustNewClient(
			c.RpcClientConf.Third,
			zrpc.WithNonBlock(),
			zrpc.WithTimeout(time.Duration(c.Timeout)*time.Millisecond),
		)
		s.SmsService = smsservice.NewSmsService(thirdClient)
		s.EmailService = emailservice.NewEmailService(thirdClient)
		s.CaptchaService = captchaservice.NewCaptchaService(thirdClient)
	}

	s.MQ = xmq.NewAsynq(s.Config.RedisConf, 1, s.Config.Log.Level)
	s.Jwt = utils.NewJwt(s.Config.Account.JwtConfig, s.Redis)
	return s
}
