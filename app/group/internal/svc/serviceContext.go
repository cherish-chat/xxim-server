package svc

import (
	"github.com/cherish-chat/xxim-server/app/group/groupmodel"
	"github.com/cherish-chat/xxim-server/app/group/internal/config"
	"github.com/cherish-chat/xxim-server/app/im/imservice"
	msgservice "github.com/cherish-chat/xxim-server/app/msg/msgService"
	"github.com/cherish-chat/xxim-server/app/relation/relationservice"
	"github.com/cherish-chat/xxim-server/app/user/userservice"
	"github.com/cherish-chat/xxim-server/common/i18n"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xconf"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config          config.Config
	zedis           *redis.Redis
	mysql           *gorm.DB
	imService       imservice.ImService
	userService     userservice.UserService
	msgService      msgservice.MsgService
	relationService relationservice.RelationService
	SystemConfigMgr *xconf.SystemConfigMgr
	*i18n.I18N
}

func NewServiceContext(c config.Config) *ServiceContext {
	s := &ServiceContext{
		Config: c,
	}
	s.SystemConfigMgr = xconf.NewSystemConfigMgr("system", c.Name, s.Mysql())
	s.I18N = i18n.NewI18N(s.Mysql())
	s.InitGroupIdCache()
	s.Mysql().AutoMigrate(
		groupmodel.Group{},
		groupmodel.GroupMember{},
	)
	return s
}

func (s *ServiceContext) Redis() *redis.Redis {
	if s.zedis == nil {
		s.zedis = s.Config.Redis.NewRedis()
	}
	return s.zedis
}

func (s *ServiceContext) Mysql() *gorm.DB {
	if s.mysql == nil {
		s.mysql = xorm.NewClient(s.Config.Mysql)
	}
	return s.mysql
}

func (s *ServiceContext) ImService() imservice.ImService {
	if s.imService == nil {
		s.imService = imservice.NewImService(zrpc.MustNewClient(s.Config.ImRpc))
	}
	return s.imService
}

func (s *ServiceContext) UserService() userservice.UserService {
	if s.userService == nil {
		s.userService = userservice.NewUserService(zrpc.MustNewClient(s.Config.UserRpc))
	}
	return s.userService
}

func (s *ServiceContext) MsgService() msgservice.MsgService {
	if s.msgService == nil {
		s.msgService = msgservice.NewMsgService(zrpc.MustNewClient(s.Config.MsgRpc))
	}
	return s.msgService
}

func (s *ServiceContext) RelationService() relationservice.RelationService {
	if s.relationService == nil {
		s.relationService = relationservice.NewRelationService(zrpc.MustNewClient(s.Config.RelationRpc))
	}
	return s.relationService
}

func (s *ServiceContext) InitGroupIdCache() {
	val, err := s.Redis().Hget(rediskey.IncrId(), rediskey.IncrGroup())
	defaultMinGroupId := s.SystemConfigMgr.GetOrDefault("default_min_group_id", "10000")
	if err != nil {
		if err == redis.Nil {
			// 不存在，从数据库中获取
			latestGroup := &groupmodel.Group{}
			err := s.Mysql().Model(latestGroup).Order("id desc").First(latestGroup).Error
			if err != nil {
				// 数据库中没有数据，从默认值开始
				val = defaultMinGroupId
			} else {
				val = latestGroup.Id
			}
			// 设置到redis中
			intGroupId := int(utils.AnyToInt64(val))
			if intGroupId < int(utils.AnyToInt64(defaultMinGroupId)) {
				intGroupId = int(utils.AnyToInt64(defaultMinGroupId))
			}
			_, err = s.Redis().Hincrby(rediskey.IncrId(), rediskey.IncrGroup(), intGroupId)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	} else {
		// 设置到redis中
		intGroupId := int(utils.AnyToInt64(val))
		if intGroupId < int(utils.AnyToInt64(defaultMinGroupId)) {
			intGroupId = int(utils.AnyToInt64(defaultMinGroupId))
			_, err = s.Redis().Hincrby(rediskey.IncrId(), rediskey.IncrGroup(), intGroupId)
			if err != nil {
				panic(err)
			}
		}
	}
}
