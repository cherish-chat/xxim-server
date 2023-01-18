package config

import (
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql                  xorm.MysqlConfig
	ImRpc                  zrpc.RpcClientConf
	NoticeRpc              zrpc.RpcClientConf
	Ip2RegionUrl           string `json:",default=https://github.com/lionsoul2014/ip2region/raw/master/data/ip2region.xdb"`
	EnableMultiDeviceLogin bool   `json:",default=true"`
}
