package config

import (
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"github.com/cherish-chat/xxim-server/common/xtdmq"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/trace"
)

type GinConfig struct {
	Port      int    `json:",default=8080"`
	Host      string `json:",default=0.0.0.0"`
	Mode      string `json:",default=dev,options=dev|pro"`
	Log       logx.LogConf
	Telemetry trace.Config `json:",optional"`
}

type Config struct {
	GinConfig
	Redis redis.RedisConf
	Mongo xmgo.Config
	TDMQ  struct {
		xtdmq.TDMQConfig
		PushConsumer     xtdmq.TDMQConsumerConfig
		InternalConsumer xtdmq.TDMQConsumerConfig
	}
}

func (c Config) SetUp() {
	logx.MustSetup(c.Log)
	trace.StartAgent(c.Telemetry)
}
