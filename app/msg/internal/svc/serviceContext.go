package svc

import (
	"github.com/cherish-chat/xxim-server/app/msg/internal/config"
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"github.com/cherish-chat/xxim-server/common/xtdmq"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config      config.Config
	msgProducer *xtdmq.TDMQProducer
	zedis       *redis.Redis
	mongo       *xmgo.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
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
