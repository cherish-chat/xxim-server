package xredis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logx"
	zeroredis "github.com/zeromicro/go-zero/core/stores/redis"
)

type Config struct {
	Host string `json:",default=localhost:6379"`
	Pass string `json:",default=123456"`
	DB   int    `json:",default=0"`
}

func GetClient(cfg Config) redis.UniversalClient {
	// 打印配置
	logx.Infof("redis config: %+v", cfg)
	opts := &redis.UniversalOptions{
		Addrs: []string{cfg.Host},
		DB:    cfg.DB,
		//PoolSize:     15,
		//MinIdleConns: 5, // redis连接池最小空闲连接数.
		Password: cfg.Pass,
		//ReadTimeout:  5,
	}
	rc := redis.NewUniversalClient(opts)
	err := rc.Ping(context.Background()).Err()
	if err != nil {
		logx.Errorf("redis ping error: %+v", err)
		panic(err)
	}
	return rc
}

func GetZeroRedis(cfg Config) *zeroredis.Redis {
	ops := make([]zeroredis.Option, 0)
	if cfg.Pass != "" {
		ops = append(ops, zeroredis.WithPass(cfg.Pass))
	}
	return zeroredis.New(cfg.Host, ops...)
}
