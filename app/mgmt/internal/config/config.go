package config

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler/middleware"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Gin GinConfig
}
type GinConfig struct {
	Cors middleware.CorsConfig
	Addr string `json:",default=0.0.0.0:6799"`
}
