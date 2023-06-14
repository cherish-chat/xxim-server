package xcache

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type xLock struct {
}

var Lock = &xLock{}

func (x *xLock) Lock(ctx context.Context, rc *redis.Redis, key string, expireSecond int) (bool, error) {
	ok, err := rc.SetnxExCtx(ctx, key, "1", expireSecond)
	if err != nil {
		return false, err
	}
	return ok, nil
}

func (x *xLock) Unlock(ctx context.Context, rc *redis.Redis, key string) error {
	_, err := rc.DelCtx(ctx, key)
	if err != nil {
		return err
	}
	return nil
}
