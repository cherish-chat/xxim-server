package config

import (
	"github.com/cherish-chat/xxim-server/common/xredis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mongo struct {
		Uri         string
		Database    string `json:",default=xxim"`
		Collections struct {
			Msg  string `json:",default=msg"`
			User string `json:",default=user"`
		}
	}
	Redis xredis.Config
}
