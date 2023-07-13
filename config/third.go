package config

import (
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
)

const ThirdService = "third.peer"

type ThirdConfig struct {
	RpcServerConf zrpc.RpcServerConf
	Mode          string
	ListenOn      string
}

func (c Config) GetThirdConfig() GatewayConfig {
	listenOn := RpcPort()
	return GatewayConfig{
		RpcServerConf: zrpc.RpcServerConf{
			ServiceConf: service.ServiceConf{
				Name:      ThirdService,
				Log:       c.GetLog(ThirdService),
				Mode:      c.Mode,
				Telemetry: c.GetJaeger(ThirdService),
			},
			ListenOn:     listenOn,
			Etcd:         c.GetEtcd(ThirdService),
			Timeout:      60000,
			CpuThreshold: 900,
		},
		Mode:     c.Mode,
		ListenOn: listenOn,
	}
}

func (c Config) GetThirdRpcConfig() zrpc.RpcClientConf {
	return zrpc.RpcClientConf{
		Etcd:     c.GetEtcd(ThirdService),
		NonBlock: true,
		Timeout:  60000,
	}
}

func (c Config) GetThirdRpcClient() (zrpc.Client, error) {
	config := c.GetThirdRpcConfig()
	return zrpc.NewClient(config)
}
