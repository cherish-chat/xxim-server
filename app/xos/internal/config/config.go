package config

import (
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type Config struct {
	Gin   GinConfig
	Mysql xorm.MysqlConfig
	Redis redis.RedisConf
	Mode  string `json:",default=debug"`
}
type GinConfig struct {
	Cors   CorsConfig
	Addr   string `json:",default=0.0.0.0:6799"`
	AesIv  string `json:",optional"`
	AesKey string `json:",optional"`
}

type CorsConfig struct {
	AllowOrigins     []string
	AllowHeaders     []string
	AllowMethods     []string
	ExposeHeaders    []string
	AllowCredentials bool `json:",default=true"`
}
