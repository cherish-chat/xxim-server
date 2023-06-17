package svc

import (
	"github.com/cherish-chat/xxim-server/app/conversation/conversationmodel"
	"github.com/cherish-chat/xxim-server/app/conversation/friendmodel"
	"github.com/cherish-chat/xxim-server/app/conversation/groupmodel"
	"github.com/cherish-chat/xxim-server/app/conversation/internal/config"
	"github.com/cherish-chat/xxim-server/app/user/client/infoservice"
	"github.com/cherish-chat/xxim-server/common/xcache"
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"github.com/qiniu/qmgo"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"time"
)

type ServiceContext struct {
	Config                        config.Config
	Redis                         *redis.Redis
	ConversationSettingCollection *qmgo.QmgoClient
	GroupCollection               *qmgo.QmgoClient
	ConversationMemberCollection  *qmgo.QmgoClient
	FriendCollection              *qmgo.QmgoClient

	InfoService infoservice.InfoService
}

func NewServiceContext(c config.Config) *ServiceContext {
	s := &ServiceContext{
		Config:                       c,
		Redis:                        xcache.MustNewRedis(c.RedisConf),
		GroupCollection:              xmgo.MustNewMongoCollection(c.Group.MongoCollection, &groupmodel.Group{}),
		ConversationMemberCollection: xmgo.MustNewMongoCollection(c.ConversationMember.MongoCollection, &conversationmodel.ConversationSetting{}),
		FriendCollection:             xmgo.MustNewMongoCollection(c.Friend.MongoCollection, &friendmodel.Friend{}),

		InfoService: infoservice.NewInfoService(zrpc.MustNewClient(
			c.RpcClientConf.User,
			zrpc.WithNonBlock(),
			zrpc.WithTimeout(time.Duration(c.Timeout)*time.Millisecond),
		)),
	}
	groupmodel.InitGroupModel(s.GroupCollection, s.Redis, s.Config.Group.MinGroupId)

	friendmodel.InitFriendModel(s.FriendCollection, s.Redis)

	conversationmodel.InitConversationMemberModel(s.ConversationMemberCollection, s.Redis)
	return s
}
