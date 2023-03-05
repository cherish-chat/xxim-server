package config

import (
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql        xorm.MysqlConfig
	ImRpc        zrpc.RpcClientConf
	UserRpc      zrpc.RpcClientConf
	MsgRpc       zrpc.RpcClientConf
	NoticeRpc    zrpc.RpcClientConf
	RelationRpc  zrpc.RpcClientConf
	Ip2RegionUrl string `json:",default=https://github.com/lionsoul2014/ip2region/raw/master/data/ip2region.xdb"`
	GroupConfig  GroupConfig
}
type GroupConfig struct {
	// 每个人能加的群数量
	MaxGroupCount int `json:",default=2000"`
	// 每个群的人数上限
	MaxGroupMemberCount int `json:",default=200000"`
}
