package config

import (
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql       xorm.MysqlConfig
	ImRpc       zrpc.RpcClientConf
	RelationRpc zrpc.RpcClientConf
	GroupRpc    zrpc.RpcClientConf
	UserRpc     zrpc.RpcClientConf
	MsgRpc      zrpc.RpcClientConf
}
