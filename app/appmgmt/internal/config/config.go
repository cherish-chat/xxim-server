package config

import (
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql        xorm.MysqlConfig
	ImRpc        zrpc.RpcClientConf
	MsgRpc       zrpc.RpcClientConf
	RelationRpc  zrpc.RpcClientConf
	UserRpc      zrpc.RpcClientConf
	GroupRpc     zrpc.RpcClientConf
	NoticeRpc    zrpc.RpcClientConf
	Ip2RegionUrl string
}
