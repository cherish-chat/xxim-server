package svc

import (
	"github.com/cherish-chat/xxim-server/app/gateway/gatewayservice"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/config"
	"github.com/zeromicro/go-zero/zrpc"
	"strings"
)

type ServiceContext struct {
	Config         config.Config
	gatewayService gatewayservice.GatewayService
}

func NewServiceContext(c config.Config) *ServiceContext {
	s := &ServiceContext{
		Config: c,
	}
	return s
}

func (s *ServiceContext) GatewayService() gatewayservice.GatewayService {
	if s.gatewayService == nil {
		listenOnSplit := strings.Split(s.Config.ListenOn, ":")
		rpcPort := listenOnSplit[len(listenOnSplit)-1]
		s.gatewayService = gatewayservice.MustNewGatewayService(zrpc.RpcClientConf{
			Endpoints: []string{"127.0.0.1:" + rpcPort},
			NonBlock:  true,
		})
	}
	return s.gatewayService
}
