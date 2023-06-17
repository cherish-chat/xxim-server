package svc

import (
	"github.com/cherish-chat/xxim-server/app/conversation/conversationmodel"
	"github.com/cherish-chat/xxim-server/app/conversation/friendmodel"
	"github.com/cherish-chat/xxim-server/app/conversation/groupmodel"
	"github.com/cherish-chat/xxim-server/app/conversation/internal/config"
	"github.com/cherish-chat/xxim-server/common/xcache"
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"github.com/qiniu/qmgo"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config                        config.Config
	Redis                         *redis.Redis
	ConversationSettingCollection *qmgo.QmgoClient
	GroupCollection               *qmgo.QmgoClient
	ConversationMemberCollection  *qmgo.QmgoClient
	FriendCollection              *qmgo.QmgoClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	s := &ServiceContext{
		Config:                       c,
		Redis:                        xcache.MustNewRedis(c.RedisConf),
		GroupCollection:              xmgo.MustNewMongoCollection(c.Group.MongoCollection, &groupmodel.Group{}),
		ConversationMemberCollection: xmgo.MustNewMongoCollection(c.ConversationMember.MongoCollection, &conversationmodel.ConversationSetting{}),
		FriendCollection:             xmgo.MustNewMongoCollection(c.Friend.MongoCollection, &friendmodel.Friend{}),
	}
	groupmodel.InitGroupModel(s.GroupCollection, s.Redis, s.Config.Group.MinGroupId)

	friendmodel.InitFriendModel(s.FriendCollection, s.Redis)

	conversationmodel.InitConversationMemberModel(s.ConversationMemberCollection, s.Redis)
	return s
}
