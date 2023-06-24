package svc

import (
	"github.com/cherish-chat/xxim-server/app/conversation/client/friendservice"
	"github.com/cherish-chat/xxim-server/app/conversation/client/groupservice"
	"github.com/cherish-chat/xxim-server/app/conversation/client/subscriptionservice"
	"github.com/cherish-chat/xxim-server/app/gateway/client/gatewayservice"
	"github.com/cherish-chat/xxim-server/app/message/client/messageservice"
	"github.com/cherish-chat/xxim-server/app/message/client/noticeservice"
	"github.com/cherish-chat/xxim-server/app/message/internal/config"
	"github.com/cherish-chat/xxim-server/app/message/messagemodel"
	"github.com/cherish-chat/xxim-server/app/message/noticemodel"
	"github.com/cherish-chat/xxim-server/app/user/client/accountservice"
	"github.com/cherish-chat/xxim-server/app/user/client/callbackservice"
	"github.com/cherish-chat/xxim-server/app/user/client/infoservice"
	"github.com/cherish-chat/xxim-server/common/xcache"
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"github.com/cherish-chat/xxim-server/common/xmq"
	"github.com/qiniu/qmgo"
	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"time"
)

type ServiceContext struct {
	Config config.Config
	Redis  *redis.Redis

	MQ xmq.MQ

	BroadcastNoticeCollection           *qmgo.QmgoClient
	SubscriptionNoticeCollection        *qmgo.QmgoClient
	SubscriptionNoticeContentCollection *qmgo.QmgoClient
	MessageCollection                   *qmgo.QmgoClient

	SendMsgTokenLimiter *limit.TokenLimiter

	//User
	CallbackService callbackservice.CallbackService
	AccountService  accountservice.AccountService
	InfoService     infoservice.InfoService
	//Conversation
	FriendService       friendservice.FriendService
	GroupService        groupservice.GroupService
	SubscriptionService subscriptionservice.SubscriptionService
	//Message
	NoticeService  noticeservice.NoticeService
	MessageService messageservice.MessageService
	//Gateway
	GatewayService gatewayservice.GatewayService
}

func NewServiceContext(c config.Config) *ServiceContext {
	userClient := zrpc.MustNewClient(
		c.RpcClientConf.User,
		zrpc.WithNonBlock(),
		zrpc.WithTimeout(time.Duration(c.Timeout)*time.Millisecond),
	)
	conversationClient := zrpc.MustNewClient(
		c.RpcClientConf.Conversation,
		zrpc.WithNonBlock(),
		zrpc.WithTimeout(time.Duration(c.Timeout)*time.Millisecond),
	)
	messageClient := zrpc.MustNewClient(
		c.RpcClientConf.Message,
		zrpc.WithNonBlock(),
		zrpc.WithTimeout(time.Duration(c.Timeout)*time.Millisecond),
	)
	gatewayClient := zrpc.MustNewClient(
		c.RpcClientConf.Gateway,
		zrpc.WithNonBlock(),
		zrpc.WithTimeout(time.Duration(c.Timeout)*time.Millisecond),
	)

	s := &ServiceContext{
		Config:                              c,
		Redis:                               xcache.MustNewRedis(c.RedisConf),
		BroadcastNoticeCollection:           xmgo.MustNewMongoCollection(c.MongoCollection.BroadcastNotice, &noticemodel.BroadcastNotice{}),
		SubscriptionNoticeCollection:        xmgo.MustNewMongoCollection(c.MongoCollection.SubscriptionNotice, &noticemodel.SubscriptionNotice{}),
		SubscriptionNoticeContentCollection: xmgo.MustNewMongoCollection(c.MongoCollection.SubscriptionNoticeContent, &noticemodel.SubscriptionNoticeContent{}),
		MessageCollection:                   xmgo.MustNewMongoCollection(c.MongoCollection.Message, &messagemodel.Message{}),

		CallbackService:     callbackservice.NewCallbackService(userClient),
		AccountService:      accountservice.NewAccountService(userClient),
		InfoService:         infoservice.NewInfoService(userClient),
		FriendService:       friendservice.NewFriendService(conversationClient),
		GroupService:        groupservice.NewGroupService(conversationClient),
		SubscriptionService: subscriptionservice.NewSubscriptionService(conversationClient),
		NoticeService:       noticeservice.NewNoticeService(messageClient),
		MessageService:      messageservice.NewMessageService(messageClient),
		GatewayService:      gatewayservice.NewGatewayService(gatewayClient),
	}

	s.SendMsgTokenLimiter = limit.NewTokenLimiter(
		c.SendMsgLimiter.Rate, c.SendMsgLimiter.Burst, s.Redis, c.SendMsgLimiter.Key,
	)

	messagemodel.InitRedisSeq(s.Redis)

	s.MQ = xmq.NewAsynq(s.Config.RedisConf, 2, s.Config.Log.Level)
	return s
}
