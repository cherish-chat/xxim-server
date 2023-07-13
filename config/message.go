package config

import (
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
)

const MessageService = "message.peer"

type MessageConfig struct {
	RpcServerConf zrpc.RpcServerConf
	Mode          string
	ListenOn      string
}

func (c Config) GetMessageConfig() GatewayConfig {
	listenOn := RpcPort()
	return GatewayConfig{
		RpcServerConf: zrpc.RpcServerConf{
			ServiceConf: service.ServiceConf{
				Name:      MessageService,
				Log:       c.GetLog(MessageService),
				Mode:      c.Mode,
				Telemetry: c.GetJaeger(MessageService),
			},
			ListenOn:     listenOn,
			Etcd:         c.GetEtcd(MessageService),
			Timeout:      60000,
			CpuThreshold: 900,
		},
		Mode:     c.Mode,
		ListenOn: listenOn,
	}
}

func (c Config) GetMessageRpcConfig() zrpc.RpcClientConf {
	return zrpc.RpcClientConf{
		Etcd:     c.GetEtcd(MessageService),
		NonBlock: true,
		Timeout:  60000,
	}
}

func (c Config) GetMessageRpcClient() (zrpc.Client, error) {
	config := c.GetMessageRpcConfig()
	return zrpc.NewClient(config)
}
