package svc

import (
	"github.com/cherish-chat/xxim-server/app/conn/connservice"
	"github.com/cherish-chat/xxim-server/app/im/internal/config"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config      config.Config
	ConnPodsMgr *connservice.ConnPodsMgr
	zedis       *redis.Redis
}

func (c *ServiceContext) Redis() *redis.Redis {
	if c.zedis == nil {
		c.zedis = c.Config.Redis.NewRedis()
	}
	return c.zedis
}

func NewServiceContext(c config.Config) *ServiceContext {
	s := &ServiceContext{
		Config: c,
	}
	s.ConnPodsMgr = connservice.NewConnPodsMgr(c.ConnRpc)
	return s
}
