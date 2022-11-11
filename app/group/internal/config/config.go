package config

import (
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mongo        xmgo.Config
	ImRpc        zrpc.RpcClientConf
	UserRpc      zrpc.RpcClientConf
	MsgRpc       zrpc.RpcClientConf
	RelationRpc  zrpc.RpcClientConf
	Ip2RegionUrl string `json:",default=https://github.com/lionsoul2014/ip2region/raw/master/data/ip2region.xdb"`
}
