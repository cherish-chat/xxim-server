package xcache

import (
	"context"
	"crypto/tls"
	goredis "github.com/go-redis/redis/v8"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"os"
)

func NewGoRedis(redisConf *redis.RedisConf) goredis.UniversalClient {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	if !redisConf.Tls {
		tlsConfig = nil
	}
	client := goredis.NewUniversalClient(&goredis.UniversalOptions{
		Addrs: []string{
			redisConf.Host,
		},
		DB:                 0,
		Dialer:             nil,
		OnConnect:          nil,
		Username:           "",
		Password:           redisConf.Pass,
		SentinelUsername:   "",
		SentinelPassword:   "",
		MaxRetries:         0,
		MinRetryBackoff:    0,
		MaxRetryBackoff:    0,
		DialTimeout:        0,
		ReadTimeout:        0,
		WriteTimeout:       0,
		PoolFIFO:           false,
		PoolSize:           0,
		MinIdleConns:       0,
		MaxConnAge:         0,
		PoolTimeout:        0,
		IdleTimeout:        0,
		IdleCheckFrequency: 0,
		TLSConfig:          tlsConfig,
		MaxRedirects:       0,
		ReadOnly:           false,
		RouteByLatency:     false,
		RouteRandomly:      false,
		MasterName:         "",
	})
	ctx, cancel := context.WithTimeout(context.Background(), redisConf.PingTimeout)
	defer cancel()
	err := client.Ping(ctx).Err()
	if err != nil {
		logx.Errorf("ping redis failed: %s", err.Error())
		os.Exit(1)
		return nil
	}
	return client
}
