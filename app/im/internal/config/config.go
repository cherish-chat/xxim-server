package config

import (
	"github.com/cherish-chat/xxim-server/common/xmq"
	"github.com/cherish-chat/xxim-server/common/xredis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	TDMQProducers struct {
		Storage xmq.TDMQProducerConfig
	}
	TDMQConsumers struct {
		Storage xmq.TDMQConsumerConfig
	}
	Mongo struct {
		Uri         string
		Database    string `json:",default=xxim"`
		Collections struct {
			Msg string `json:",default=msg"`
		}
	}
	Redis xredis.Config
}
