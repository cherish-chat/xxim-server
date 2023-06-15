package svc

import (
	"github.com/cherish-chat/xxim-server/app/third/thirdservice"
	"github.com/cherish-chat/xxim-server/app/user/internal/config"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xcache"
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"github.com/cherish-chat/xxim-server/common/xmq"
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
	MQ             xmq.MQ
	Jwt            *utils.Jwt
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
	s.MQ = xmq.NewAsynq(s.Config.RedisConf, 1, s.Config.Log.Level)
	s.Jwt = utils.NewJwt(s.Config.Account.JwtConfig, s.Redis)
	return s
}
