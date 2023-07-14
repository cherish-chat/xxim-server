package svc

import (
	"github.com/cherish-chat/xxim-server/app/api/gateway/client/internalservice"
	"github.com/cherish-chat/xxim-server/app/service/conversation/client/channelservice"
	"github.com/cherish-chat/xxim-server/app/service/conversation/client/groupservice"
	"github.com/cherish-chat/xxim-server/app/service/message/messagemodel"
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
	MQ     xmq.MQ

	MessageCollection *qmgo.QmgoClient
	NoticeCollection  *qmgo.QmgoClient

	GroupService   groupservice.GroupService
	ChannelService channelservice.ChannelService

	InternalService internalservice.InternalService
}

func NewServiceContext(c config.Config) *ServiceContext {
	conversationClient, err := c.GetConversationRpcClient()
	if err != nil {
		logx.Errorf("get conversation rpc client error: %s", err.Error())
		panic(err)
	}
	gatewayClient, err := c.GetGatewayRpcClient()
	if err != nil {
		logx.Errorf("get gateway rpc client error: %s", err.Error())
		panic(err)
	}
	s := &ServiceContext{
		Config: c,
		Redis:  c.GetRedis(3),

		MessageCollection: xmgo.MustNewMongoCollection(c.MongoCollection.Message, &messagemodel.Message{}),
		NoticeCollection:  xmgo.MustNewMongoCollection(c.MongoCollection.Notice, &messagemodel.Message{}),

		GroupService:   groupservice.NewGroupService(conversationClient),
		ChannelService: channelservice.NewChannelService(conversationClient),

		InternalService: internalservice.NewInternalService(gatewayClient),
	}
	messagemodel.InitRedisSeq(c.GetRedis(3))
	s.MQ = xmq.NewAsynq(c.GetZeroRedisConf(), 3, s.Config.Log.Level)
	return s
}
