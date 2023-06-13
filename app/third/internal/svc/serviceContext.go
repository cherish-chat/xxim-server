package svc

import (
	"github.com/cherish-chat/xxim-server/app/third/internal/config"
	"github.com/cherish-chat/xxim-server/common/xcache"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	Redis  *redis.Redis
	Mysql  *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
		Redis:  xcache.MustNewRedis(c.RedisConf),
		Mysql:  xorm.MustNewMysql(c.MysqlConf),
	}
}
