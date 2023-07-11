package config

import (
	"errors"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
)

const GatewayService = "gateway.peer"

type GatewayConfig struct {
	RpcServerConf zrpc.RpcServerConf
	Mode          string
	ListenOn      string
}

func (c Config) GetGatewayRpcConfig() zrpc.RpcClientConf {
	return zrpc.RpcClientConf{
		Etcd:     c.GetEtcd(GatewayService),
		NonBlock: true,
		Timeout:  60000,
	}
}

func (c Config) GetGatewayRpcClient() (zrpc.Client, error) {
	config := c.GetGatewayRpcConfig()
	return zrpc.NewClient(config)
}

func (c Config) GetAllGatewayRpcClient() ([]zrpc.Client, error) {
	clients, ok := zrpc.GetServiceClients(GatewayService)
	if !ok {
		return nil, errors.New("no gateway.peer client found")
	}
	return clients, nil
}

func (c Config) GetGatewayConfig() GatewayConfig {
	listenOn := RpcPort()
	return GatewayConfig{
		RpcServerConf: zrpc.RpcServerConf{
			ServiceConf: service.ServiceConf{
				Name:      GatewayService,
				Log:       c.GetLog(GatewayService),
				Mode:      c.Mode,
				Telemetry: c.GetJaeger(GatewayService),
			},
			ListenOn:     listenOn,
			Etcd:         c.GetEtcd(GatewayService),
			Timeout:      60000,
			CpuThreshold: 900,
		},
		Mode:     c.Mode,
		ListenOn: listenOn,
	}
}
