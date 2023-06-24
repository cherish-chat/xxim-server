package svc

import (
	"github.com/cherish-chat/xxim-server/app/conversation/client/conversationservice"
	"github.com/cherish-chat/xxim-server/app/conversation/client/subscriptionservice"
	"github.com/cherish-chat/xxim-server/app/conversation/conversationmodel"
	"github.com/cherish-chat/xxim-server/app/conversation/friendmodel"
	"github.com/cherish-chat/xxim-server/app/conversation/groupmodel"
	"github.com/cherish-chat/xxim-server/app/conversation/internal/config"
	"github.com/cherish-chat/xxim-server/app/conversation/subscriptionmodel"
	"github.com/cherish-chat/xxim-server/app/message/client/messageservice"
	"github.com/cherish-chat/xxim-server/app/message/client/noticeservice"
	"github.com/cherish-chat/xxim-server/app/user/client/infoservice"
	"github.com/cherish-chat/xxim-server/common/xcache"
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"github.com/qiniu/qmgo"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"time"
)

type ServiceContext struct {
	Config                          config.Config
	Redis                           *redis.Redis
	ConversationSettingCollection   *qmgo.QmgoClient
	GroupCollection                 *qmgo.QmgoClient
	GroupSubscribeCollection        *qmgo.QmgoClient
	ConversationMemberCollection    *qmgo.QmgoClient
	FriendCollection                *qmgo.QmgoClient
	FriendApplyRecordCollection     *qmgo.QmgoClient
	SubscriptionCollection          *qmgo.QmgoClient
	UserSubscriptionCollection      *qmgo.QmgoClient
	SubscriptionSubscribeCollection *qmgo.QmgoClient

	InfoService         infoservice.InfoService
	NoticeService       noticeservice.NoticeService
	MessageService      messageservice.MessageService
	SubscriptionService subscriptionservice.SubscriptionService
	ConversationService conversationservice.ConversationService
}

func NewServiceContext(c config.Config) *ServiceContext {
	userClient := zrpc.MustNewClient(
		c.RpcClientConf.User,
		zrpc.WithNonBlock(),
		zrpc.WithTimeout(time.Duration(c.Timeout)*time.Millisecond),
	)
	messageClient := zrpc.MustNewClient(
		c.RpcClientConf.Message,
		zrpc.WithNonBlock(),
		zrpc.WithTimeout(time.Duration(c.Timeout)*time.Millisecond),
	)
	conversationClient := zrpc.MustNewClient(
		c.RpcClientConf.Conversation,
		zrpc.WithNonBlock(),
		zrpc.WithTimeout(time.Duration(c.Timeout)*time.Millisecond),
	)

	s := &ServiceContext{
		Config:                          c,
		Redis:                           xcache.MustNewRedis(c.RedisConf),
		GroupCollection:                 xmgo.MustNewMongoCollection(c.MongoCollection.Group, &groupmodel.Group{}),
		GroupSubscribeCollection:        xmgo.MustNewMongoCollection(c.MongoCollection.GroupSubscribe, &groupmodel.GroupSubscribe{}),
		ConversationMemberCollection:    xmgo.MustNewMongoCollection(c.MongoCollection.ConversationMember, &conversationmodel.ConversationMember{}),
		FriendCollection:                xmgo.MustNewMongoCollection(c.MongoCollection.Friend, &friendmodel.Friend{}),
		FriendApplyRecordCollection:     xmgo.MustNewMongoCollection(c.MongoCollection.FriendApplyRecord, &friendmodel.FriendApplyRecord{}),
		SubscriptionCollection:          xmgo.MustNewMongoCollection(c.MongoCollection.Subscription, &subscriptionmodel.Subscription{}),
		UserSubscriptionCollection:      xmgo.MustNewMongoCollection(c.MongoCollection.UserSubscription, &subscriptionmodel.UserSubscription{}),
		SubscriptionSubscribeCollection: xmgo.MustNewMongoCollection(c.MongoCollection.SubscriptionSubscribe, &subscriptionmodel.SubscriptionSubscribe{}),

		InfoService:         infoservice.NewInfoService(userClient),
		NoticeService:       noticeservice.NewNoticeService(messageClient),
		MessageService:      messageservice.NewMessageService(messageClient),
		SubscriptionService: subscriptionservice.NewSubscriptionService(conversationClient),
		ConversationService: conversationservice.NewConversationService(conversationClient),
	}
	groupmodel.InitGroupModel(s.GroupCollection, s.Redis, s.Config.Group.MinGroupId)

	friendmodel.InitFriendModel(s.FriendCollection, s.Redis)

	conversationmodel.InitConversationMemberModel(s.ConversationMemberCollection, s.Redis)

	subscriptionmodel.InitSystemSubscription(s.SubscriptionCollection)
	return s
}
