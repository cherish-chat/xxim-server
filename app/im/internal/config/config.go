package config

import (
	"github.com/cherish-chat/xxim-server/app/conn/connservice"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	ConnRpc      connservice.ConnPodsConfig
	Mysql        xorm.MysqlConfig
	Ip2RegionUrl string `json:",default=https://github.com/lionsoul2014/ip2region/raw/master/data/ip2region.xdb"`
	MsgRpc       zrpc.RpcClientConf
	NoticeRpc    zrpc.RpcClientConf
}
