package svc

import (
	"github.com/cherish-chat/xxim-server/app/group/groupservice"
	"github.com/cherish-chat/xxim-server/app/im/imservice"
	"github.com/cherish-chat/xxim-server/app/msg/internal/config"
	"github.com/cherish-chat/xxim-server/app/relation/relationservice"
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"github.com/cherish-chat/xxim-server/common/xtdmq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config          config.Config
	msgProducer     *xtdmq.TDMQProducer
	zedis           *redis.Redis
	mongo           *xmgo.Client
	imService       imservice.ImService
	relationService relationservice.RelationService
	groupService    groupservice.GroupService
}

func NewServiceContext(c config.Config) *ServiceContext {
	s := &ServiceContext{
		Config: c,
	}
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

func (s *ServiceContext) Mongo() *xmgo.Client {
	if s.mongo == nil {
		s.mongo = xmgo.NewClient(s.Config.Mongo)
	}
	return s.mongo
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
