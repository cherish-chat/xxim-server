package config

import (
	"github.com/cherish-chat/xxim-server/app/conn/connservice"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	ConnRpc connservice.ConnPodsConfig
}
