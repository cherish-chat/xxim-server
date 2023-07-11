package config

import (
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
)

const UserService = "user.peer"

type UserConfig struct {
	RpcServerConf zrpc.RpcServerConf
	Mode          string
	ListenOn      string
}

func (c Config) GetUserConfig() GatewayConfig {
	listenOn := RpcPort()
	return GatewayConfig{
		RpcServerConf: zrpc.RpcServerConf{
			ServiceConf: service.ServiceConf{
				Name:      UserService,
				Log:       c.GetLog(UserService),
				Mode:      c.Mode,
				Telemetry: c.GetJaeger(UserService),
			},
			ListenOn:     listenOn,
			Etcd:         c.GetEtcd(UserService),
			Timeout:      60000,
			CpuThreshold: 900,
		},
		Mode:     c.Mode,
		ListenOn: listenOn,
	}
}

func (c Config) GetUserRpcConfig() zrpc.RpcClientConf {
	return zrpc.RpcClientConf{
		Etcd:     c.GetEtcd(UserService),
		NonBlock: true,
		Timeout:  60000,
	}
}

func (c Config) GetUserRpcClient() (zrpc.Client, error) {
	config := c.GetUserRpcConfig()
	return zrpc.NewClient(config)
}
