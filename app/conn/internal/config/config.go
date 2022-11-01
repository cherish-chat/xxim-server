package config

import (
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type HttpConfig struct {
	Port int    `json:",default=8080"`
	Host string `json:",default=0.0.0.0"`
}

type Config struct {
	zrpc.RpcServerConf
	Websocket HttpConfig
	Redis     redis.RedisConf
	Mongo     xmgo.Config
	ImRpc     zrpc.RpcClientConf
}
