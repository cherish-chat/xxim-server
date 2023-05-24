package config

import "github.com/zeromicro/go-zero/zrpc"

type Config struct {
	zrpc.RpcServerConf
	Http struct {
		Encrypt struct {
			Enable bool   `json:",optional"`
			AesKey string `json:",optional"`
			AesIv  string `json:",optional"`
		} `json:",optional"`
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
}
