package svc

import (
	"github.com/cherish-chat/xxim-server/app/conn/internal/config"
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config config.Config
	Redis  *redis.Redis
	Mgo    *xmgo.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	s := &ServiceContext{
		Config: c,
	}
	s.mustSetup()
	return s
}

func (s *ServiceContext) mustSetup() {
	s.Config.SetUp()
	s.Redis = s.Config.Redis.NewRedis()
	s.Mgo = xmgo.NewClient(s.Config.Mongo)
}
