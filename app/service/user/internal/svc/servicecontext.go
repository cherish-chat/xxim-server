package svc

import (
	"github.com/cherish-chat/xxim-server/app/service/conversation/client/channelservice"
	"github.com/cherish-chat/xxim-server/app/service/conversation/client/friendservice"
	"github.com/cherish-chat/xxim-server/app/service/conversation/client/groupservice"
	"github.com/cherish-chat/xxim-server/app/service/third/client/captchaservice"
	"github.com/cherish-chat/xxim-server/app/service/third/client/emailservice"
	"github.com/cherish-chat/xxim-server/app/service/third/client/smsservice"
	"github.com/cherish-chat/xxim-server/app/service/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"github.com/cherish-chat/xxim-server/common/xmq"
	"github.com/cherish-chat/xxim-server/config"
	"github.com/go-redis/redis/v8"
	"github.com/qiniu/qmgo"
	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	Config config.Config
	Redis  redis.UniversalClient

	UserCollection        *qmgo.QmgoClient
	UserSettingCollection *qmgo.QmgoClient

	MQ  xmq.MQ
	Jwt *utils.Jwt

	SmsService     smsservice.SmsService
	EmailService   emailservice.EmailService
	CaptchaService captchaservice.CaptchaService
	FriendService  friendservice.FriendService
	GroupService   groupservice.GroupService
	ChannelService channelservice.ChannelService
}

func NewServiceContext(c config.Config) *ServiceContext {
	s := &ServiceContext{
		Config:                c,
		Redis:                 c.GetRedis(1),
		UserCollection:        xmgo.MustNewMongoCollection(c.MongoCollection.User, &usermodel.User{}),
		UserSettingCollection: xmgo.MustNewMongoCollection(c.MongoCollection.UserSetting, &usermodel.UserSetting{}),
	}

	//third rpc
	{
		thirdClient, err := c.GetThirdRpcClient()
		if err != nil {
			logx.Errorf("get third rpc client error: %v", err)
			panic(err)
		}
		s.SmsService = smsservice.NewSmsService(thirdClient)
		s.EmailService = emailservice.NewEmailService(thirdClient)
		s.CaptchaService = captchaservice.NewCaptchaService(thirdClient)
	}
	// conversation rpc
	{
		conversationClient, err := c.GetConversationRpcClient()
		if err != nil {
			logx.Errorf("get conversation rpc client error: %v", err)
			panic(err)
		}
		s.FriendService = friendservice.NewFriendService(conversationClient)
		s.GroupService = groupservice.NewGroupService(conversationClient)
		s.ChannelService = channelservice.NewChannelService(conversationClient)
	}

	s.MQ = xmq.NewAsynq(s.Config.GetZeroRedisConf(), 1, s.Config.Log.Level)
	s.Jwt = utils.NewJwt(s.Config.Account.Login.JwtConfig, s.Redis)
	return s
}
