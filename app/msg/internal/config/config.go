package config

import (
	"github.com/cherish-chat/xxim-server/app/conn/connservice"
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"github.com/cherish-chat/xxim-server/common/xtdmq"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	TDMQ struct {
		xtdmq.TDMQConfig
		xtdmq.TDMQConsumerConfig
		Producers struct {
			Msg xtdmq.TDMQProducerConfig
		}
	}
	Mongo   xmgo.Config
	ImRpc   zrpc.RpcClientConf
	ConnRpc connservice.ConnPodsConfig
}
