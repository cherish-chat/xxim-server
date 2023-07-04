package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Http struct {
		Cors struct {
			Enable           bool     `json:",optional"`
			AllowOrigins     []string `json:",optional"`
			AllowHeaders     []string `json:",optional"`
			AllowMethods     []string `json:",optional"`
			ExposeHeaders    []string `json:",optional"`
			AllowCredentials bool     `json:",optional"`
		} `json:",optional"`
		ApiLog struct {
			Apis []string `json:",optional"` // 格式: GET r'^/api/v1/user/.*' 表示所有以 /api/v1/user/ 开头的 GET 请求都会被记录
		}
		Host string `json:",default=0.0.0.0"`
		Port int    `json:",default=34500"`
	}
	Cloudx struct {
		StunUrls        []string
		Host            string
		Port            int
		Ssl             bool `json:",default=false"`
		AppId           string
		ClientId        string
		ClientSecret    string
		KeepLiveSeconds int `json:",default=30"`
	} `json:",optional"`
	Websocket struct {
		KeepAliveTickerSecond int `json:",default=30"` // 定时器，每隔n秒检测连接是否存活
		KeepAliveSecond       int `json:",default=60"` // 检测是否存活时，如果超过n秒没有收到客户端的消息，则关闭连接
	}
	RpcClientConf struct {
		Dispatch     zrpc.RpcClientConf
		User         zrpc.RpcClientConf
		Conversation zrpc.RpcClientConf
		Third        zrpc.RpcClientConf
		Message      zrpc.RpcClientConf
	}
	RedisConf redis.RedisConf
}
