package config

import (
	"github.com/zeromicro/go-zero/zrpc"
)

type HttpConfig struct {
	Port int    `json:",default=6701"`
	Host string `json:",default=0.0.0.0"`
}

type Config struct {
	zrpc.RpcServerConf
	Websocket     HttpConfig
	ImRpc         zrpc.RpcClientConf
	MsgRpc        zrpc.RpcClientConf
	RelationRpc   zrpc.RpcClientConf
	UserRpc       zrpc.RpcClientConf
	GroupRpc      zrpc.RpcClientConf
	AppMgmtRpc    zrpc.RpcClientConf
	NoticeRpc     zrpc.RpcClientConf
	RsaPublicKey  string // 客户端使用公钥来加密
	RsaPrivateKey string // 服务端使用私钥来解密
	Ip2RegionUrl  string
}
