package config

import (
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
)

const ConversationService = "conversation.peer"

type ConversationConfig struct {
	RpcServerConf zrpc.RpcServerConf
	Mode          string
	ListenOn      string
}

func (c Config) GetConversationConfig() GatewayConfig {
	listenOn := RpcPort()
	return GatewayConfig{
		RpcServerConf: zrpc.RpcServerConf{
			ServiceConf: service.ServiceConf{
				Name:      ConversationService,
				Log:       c.GetLog(ConversationService),
				Mode:      c.Mode,
				Telemetry: c.GetJaeger(ConversationService),
			},
			ListenOn:     listenOn,
			Etcd:         c.GetEtcd(ConversationService),
			Timeout:      60000,
			CpuThreshold: 900,
		},
		Mode:     c.Mode,
		ListenOn: listenOn,
	}
}

func (c Config) GetConversationRpcConfig() zrpc.RpcClientConf {
	return zrpc.RpcClientConf{
		Etcd:     c.GetEtcd(ConversationService),
		NonBlock: true,
		Timeout:  60000,
	}
}

func (c Config) GetConversationRpcClient() (zrpc.Client, error) {
	config := c.GetConversationRpcConfig()
	return zrpc.NewClient(config)
}
