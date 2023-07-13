package xcache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type xLock struct {
}

var Lock = &xLock{}

func (x *xLock) Lock(ctx context.Context, rc redis.UniversalClient, key string, expireSecond int) (bool, error) {
	ok, err := rc.SetNX(ctx, key, "1", time.Duration(expireSecond)*time.Second).Result()
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (x *xLock) Unlock(ctx context.Context, rc redis.UniversalClient, key string) error {
	_, err := rc.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
}
