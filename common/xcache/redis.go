package xcache

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"os"
)

func MustNewRedis(config redis.RedisConf) *redis.Redis {
	client, err := redis.NewRedis(config)
	if err != nil {
		logx.Errorf("redis.NewRedis error: %v", err)
		os.Exit(1)
		return nil
	}
	return client
}
