package xredis

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// HMSetEx sets the hash with ttl.
func HMSetEx(rc *redis.Redis, ctx context.Context, key string, fields map[string]string, ttl int) error {
	var err error
	batch := 1000
	for i := 0; i < len(fields); i += batch {
		var fields2 = make(map[string]string, batch)
		for k, v := range fields {
			fields2[k] = v
		}
		err = rc.HmsetCtx(ctx, key, fields2)
		if err != nil {
			return err
		}
	}
	err = rc.Expire(key, ttl)
	return err
}
