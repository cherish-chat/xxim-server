package svc

import (
	"github.com/cherish-chat/xxim-server/app/third/thirdservice"
	"github.com/cherish-chat/xxim-server/app/user/internal/config"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/xcache"
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"github.com/qiniu/qmgo"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"time"
)

type ServiceContext struct {
	Config         config.Config
	Redis          *redis.Redis
	UserCollection *qmgo.QmgoClient
	ThirdService   thirdservice.ThirdService
}

func NewServiceContext(c config.Config) *ServiceContext {
	s := &ServiceContext{
		Config:         c,
		Redis:          xcache.MustNewRedis(c.RedisConf),
		UserCollection: xmgo.MustNewMongoCollection(c.User.MongoCollection, &usermodel.User{}),
	}
	s.ThirdService = thirdservice.NewThirdService(zrpc.MustNewClient(
		c.RpcClientConf.Third,
		zrpc.WithNonBlock(),
		zrpc.WithTimeout(time.Duration(c.Timeout)*time.Millisecond),
	))
	return s
}
