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
	GroupMemberCollection         *qmgo.QmgoClient
	FriendCollection              *qmgo.QmgoClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	s := &ServiceContext{
		Config:                        c,
		Redis:                         xcache.MustNewRedis(c.RedisConf),
		ConversationSettingCollection: xmgo.MustNewMongoCollection(c.ConversationSetting.MongoCollection, &conversationmodel.ConversationSetting{}),
		GroupCollection:               xmgo.MustNewMongoCollection(c.Group.MongoCollection, &groupmodel.Group{}),
		GroupMemberCollection:         xmgo.MustNewMongoCollection(c.GroupMember.MongoCollection, &groupmodel.GroupMember{}),
		FriendCollection:              xmgo.MustNewMongoCollection(c.Friend.MongoCollection, &friendmodel.Friend{}),
	}
	groupmodel.InitGroupModel(s.GroupCollection, s.Redis, s.Config.Group.MinGroupId)
	groupmodel.InitGroupMemberModel(s.GroupMemberCollection, s.Redis)

	friendmodel.InitFriendModel(s.FriendCollection, s.Redis)

	conversationmodel.InitConversationSettingModel(s.ConversationSettingCollection, s.Redis)
	return s
}
