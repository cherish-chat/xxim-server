package svc

import (
	"github.com/cherish-chat/xxim-server/app/user/internal/config"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/xcache"
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"github.com/qiniu/qmgo"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config         config.Config
	Redis          *redis.Redis
	UserCollection *qmgo.QmgoClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	s := &ServiceContext{
		Config:         c,
		Redis:          xcache.MustNewRedis(c.RedisConf),
		UserCollection: xmgo.MustNewMongoCollection(c.User.MongoCollection, &usermodel.User{}),
	}
	usermodel.InitUserModel(s.UserCollection, s.Redis)
	return s
}
