package config

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler/middleware"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Gin         GinConfig
	ImRpc       zrpc.RpcClientConf
	MsgRpc      zrpc.RpcClientConf
	RelationRpc zrpc.RpcClientConf
	UserRpc     zrpc.RpcClientConf
	GroupRpc    zrpc.RpcClientConf
	NoticeRpc   zrpc.RpcClientConf
}
type GinConfig struct {
	Cors middleware.CorsConfig
	Addr string `json:",default=0.0.0.0:6799"`
}
