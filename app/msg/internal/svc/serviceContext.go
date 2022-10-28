package svc

import (
	"github.com/cherish-chat/xxim-server/app/msg/internal/config"
	"github.com/cherish-chat/xxim-server/common/xtdmq"
)

type ServiceContext struct {
	Config      config.Config
	msgProducer *xtdmq.TDMQProducer
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
