package svc

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/config"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config config.Config
	redis  *redis.Redis
}

func NewServiceContext(c config.Config, rc *redis.Redis) *ServiceContext {
	return &ServiceContext{
		Config: c,
		redis:  rc,
	}
}

func (s *ServiceContext) Redis() *redis.Redis {
	return s.redis
}
