package config

import (
	"github.com/cherish-chat/xxim-server/app/conn/connservice"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	Redis        redis.RedisKeyConf
	UserRpc      zrpc.RpcClientConf
	RelationRpc  zrpc.RpcClientConf
	GroupRpc     zrpc.RpcClientConf
	MsgRpc       zrpc.RpcClientConf
	ImRpc        zrpc.RpcClientConf
	NoticeRpc    zrpc.RpcClientConf
	Ip2RegionUrl string `json:",default=https://github.com/lionsoul2014/ip2region/raw/master/data/ip2region.xdb"`
	EnableSSL    bool   `json:",default=false"`
	ConnRpc      connservice.ConnPodsConfig
	Mysql        xorm.MysqlConfig
}
