package svc

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtservice"
	"github.com/cherish-chat/xxim-server/app/xos/internal/config"
	"github.com/cherish-chat/xxim-server/common/xconf"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config      config.Config
	redis       *redis.Redis
	mysql       *gorm.DB
	ConfigMgr   *xconf.ConfigMgr
	mgmtService mgmtservice.MgmtService
}

func NewServiceContext(config config.Config, mgmtService mgmtservice.MgmtService) *ServiceContext {
	s := &ServiceContext{Config: config, mgmtService: mgmtService}
	s.ConfigMgr = xconf.NewConfigMgr(s.Mysql(), s.Redis(), "system")
	return s
}

func (s *ServiceContext) Redis() *redis.Redis {
	if s.redis == nil {
		s.redis = s.Config.Redis.NewRedis()
	}
	return s.redis
}

func (s *ServiceContext) Mysql() *gorm.DB {
	if s.mysql == nil {
		s.mysql = xorm.NewClient(s.Config.Mysql)
	}
	return s.mysql
}

func (s *ServiceContext) MgmtService() mgmtservice.MgmtService {
	return s.mgmtService
}
