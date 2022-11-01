package svc

import (
	"github.com/cherish-chat/xxim-server/app/conn/connservice"
	"github.com/cherish-chat/xxim-server/app/im/internal/config"
)

type ServiceContext struct {
	Config      config.Config
	ConnPodsMgr *connservice.ConnPodsMgr
}

func NewServiceContext(c config.Config) *ServiceContext {
	s := &ServiceContext{
		Config: c,
	}
	s.ConnPodsMgr = connservice.NewConnPodsMgr(c.ConnRpc)
	return s
}
