package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type HttpConfig struct {
	Port int    `json:",default=8080"`
	Host string `json:",default=0.0.0.0"`
}

type Config struct {
	zrpc.RpcServerConf
	Websocket HttpConfig
	ImRpc     zrpc.RpcClientConf
	MsgRpc    zrpc.RpcClientConf
	NoticeRpc zrpc.RpcClientConf
}
