package config

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler/middleware"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql       xorm.MysqlConfig
	Gin         GinConfig
	ImRpc       zrpc.RpcClientConf
	MsgRpc      zrpc.RpcClientConf
	RelationRpc zrpc.RpcClientConf
	UserRpc     zrpc.RpcClientConf
	GroupRpc    zrpc.RpcClientConf
	NoticeRpc   zrpc.RpcClientConf
	SuperAdmin  struct {
		Id       string `json:",default=superadmin"`
		Password string `json:",default=superadmin"` // 只有该管理员未创建时才会创建并设置密码 后续不会修改密码
	}
	Ip2RegionUrl string
}
type GinConfig struct {
	Cors middleware.CorsConfig
	Addr string `json:",default=0.0.0.0:6799"`
}
