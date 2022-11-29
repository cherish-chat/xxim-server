package svc

import (
	"github.com/cherish-chat/xxim-server/app/group/groupservice"
	"github.com/cherish-chat/xxim-server/app/im/imservice"
	"github.com/cherish-chat/xxim-server/app/msg/internal/config"
	"github.com/cherish-chat/xxim-server/app/msg/msgmodel"
	"github.com/cherish-chat/xxim-server/app/relation/relationservice"
	"github.com/cherish-chat/xxim-server/common/pkg/mobpush"
	"github.com/cherish-chat/xxim-server/common/xconf"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xtdmq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config          config.Config
	msgProducer     *xtdmq.TDMQProducer
	zedis           *redis.Redis
	mysql           *gorm.DB
	imService       imservice.ImService
	relationService relationservice.RelationService
	groupService    groupservice.GroupService
	MobPush         *mobpush.Pusher
	SystemConfigMgr *xconf.SystemConfigMgr
}

func NewServiceContext(c config.Config) *ServiceContext {
	s := &ServiceContext{
		Config: c,
	}
	s.Mysql().AutoMigrate(
		msgmodel.Msg{},
		xorm.HashKv{},
	)
	s.MobPush = mobpush.NewPusher(c.MobPush)
	s.SystemConfigMgr = xconf.NewSystemConfigMgr("system", c.Name, s.Mysql())
	return s
}

func (s *ServiceContext) MsgProducer() *xtdmq.TDMQProducer {
	if s.msgProducer == nil {
		s.msgProducer = xtdmq.NewTDMQProducer(s.Config.TDMQ.TDMQConfig, s.Config.TDMQ.Producers.Msg)
	}
	return s.msgProducer
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

func (s *ServiceContext) RelationService() relationservice.RelationService {
	if s.relationService == nil {
		s.relationService = relationservice.NewRelationService(zrpc.MustNewClient(s.Config.RelationRpc))
	}
	return s.relationService
}

func (s *ServiceContext) GroupService() groupservice.GroupService {
	if s.groupService == nil {
		s.groupService = groupservice.NewGroupService(zrpc.MustNewClient(s.Config.GroupRpc))
	}
	return s.groupService
}
