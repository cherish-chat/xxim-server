package svc

import (
	"github.com/cherish-chat/xxim-server/app/third/internal/config"
	"github.com/cherish-chat/xxim-server/common/xcache"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config config.Config
	Redis  *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Redis:  xcache.MustNewRedis(c.RedisConf),
	}
}
