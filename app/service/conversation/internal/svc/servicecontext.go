package svc

import (
	"github.com/cherish-chat/xxim-server/app/service/conversation/channelmodel"
	"github.com/cherish-chat/xxim-server/app/service/conversation/client/channelservice"
	"github.com/cherish-chat/xxim-server/app/service/conversation/client/sessionservice"
	"github.com/cherish-chat/xxim-server/app/service/conversation/friendmodel"
	"github.com/cherish-chat/xxim-server/app/service/conversation/groupmodel"
	"github.com/cherish-chat/xxim-server/app/service/message/client/messageservice"
	"github.com/cherish-chat/xxim-server/app/service/message/client/noticeservice"
	"github.com/cherish-chat/xxim-server/app/service/user/client/userservice"
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"github.com/cherish-chat/xxim-server/config"
	"github.com/go-redis/redis/v8"
	"github.com/qiniu/qmgo"
	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	Config config.Config
	Redis  redis.UniversalClient

	ConversationSettingCollection   *qmgo.QmgoClient
	GroupCollection                 *qmgo.QmgoClient
	GroupSubscribeCacheCollection   *qmgo.QmgoClient
	GroupMemberCollection           *qmgo.QmgoClient
	FriendCollection                *qmgo.QmgoClient
	FriendApplyRecordCollection     *qmgo.QmgoClient
	ChannelCollection               *qmgo.QmgoClient
	ChannelSubscribeCacheCollection *qmgo.QmgoClient
	ChannelMemberCollection         *qmgo.QmgoClient

	UserService userservice.UserService

	NoticeService  noticeservice.NoticeService
	MessageService messageservice.MessageService

	ChannelService channelservice.ChannelService
	SessionService sessionservice.SessionService
}

func NewServiceContext(c config.Config) *ServiceContext {
	userClient, err := c.GetUserRpcClient()
	if err != nil {
		logx.Errorf("get user rpc client error: %s", err.Error())
		panic(err)
	}
	conversationClient, err := c.GetConversationRpcClient()
	if err != nil {
		logx.Errorf("get conversation rpc client error: %s", err.Error())
		panic(err)
	}
	messageClient, err := c.GetMessageRpcClient()
	if err != nil {
		logx.Errorf("get message rpc client error: %s", err.Error())
		panic(err)
	}
	s := &ServiceContext{
		Config:                        c,
		Redis:                         c.GetRedis(2),
		GroupCollection:               xmgo.MustNewMongoCollection(c.MongoCollection.Group, &groupmodel.Group{}),
		GroupSubscribeCacheCollection: xmgo.MustNewMongoCollection(c.MongoCollection.GroupSubscribeCache, &groupmodel.GroupSubscribeCache{}),
		GroupMemberCollection:         xmgo.MustNewMongoCollection(c.MongoCollection.GroupMember, &groupmodel.GroupMember{}),

		FriendCollection:            xmgo.MustNewMongoCollection(c.MongoCollection.Friend, &friendmodel.Friend{}),
		FriendApplyRecordCollection: xmgo.MustNewMongoCollection(c.MongoCollection.FriendApplyRecord, &friendmodel.FriendApplyRecord{}),

		ChannelCollection:               xmgo.MustNewMongoCollection(c.MongoCollection.Channel, &channelmodel.Channel{}),
		ChannelSubscribeCacheCollection: xmgo.MustNewMongoCollection(c.MongoCollection.ChannelSubscribeCache, &channelmodel.ChannelSubscribeCache{}),
		ChannelMemberCollection:         xmgo.MustNewMongoCollection(c.MongoCollection.ChannelMember, &channelmodel.ChannelMember{}),

		UserService:    userservice.NewUserService(userClient),
		NoticeService:  noticeservice.NewNoticeService(messageClient),
		MessageService: messageservice.NewMessageService(messageClient),
		ChannelService: channelservice.NewChannelService(conversationClient),
		SessionService: sessionservice.NewSessionService(conversationClient),
	}

	return s
}
