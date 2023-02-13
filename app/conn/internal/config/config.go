package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type HttpConfig struct {
	Port int    `json:",default=6701"`
	Host string `json:",default=0.0.0.0"`
}

type Config struct {
	zrpc.RpcServerConf
	Websocket   HttpConfig
	ImRpc       zrpc.RpcClientConf
	MsgRpc      zrpc.RpcClientConf
	RelationRpc zrpc.RpcClientConf
	UserRpc     zrpc.RpcClientConf
	GroupRpc    zrpc.RpcClientConf
	AppMgmtRpc  zrpc.RpcClientConf
	NoticeRpc   zrpc.RpcClientConf
}
