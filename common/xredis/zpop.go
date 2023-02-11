package xredis

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

const (
	// zpopmin
	ZPopMinScript = `
local key = KEYS[1]
local limit = 1
local result = redis.call('zrange', key, 0, limit - 1)
if #result > 0 then
	redis.call('zrem', key, result[1])
	return result[1]
else
	return nil
end
`
	// zpopmax
	ZPopMaxScript = `
local key = KEYS[1]
local limit = 1
local result = redis.call('zrevrange', key, 0, limit - 1)
if #result > 0 then
	redis.call('zrem', key, result[1])
	return result[1]
else
	return nil
end
`
)

var (
	ZPopMinSha1 string
	ZPopMaxSha1 string
)

func ZPopMin(ctx context.Context, rc *redis.Redis, key string) (string, error) {
	if ZPopMinSha1 == "" {
		var err error
		ZPopMinSha1, err = rc.ScriptLoadCtx(ctx, ZPopMinScript)
		if err != nil {
			return "", err
		}
	}
	val, err := rc.EvalShaCtx(ctx, ZPopMinSha1, []string{key})
	if err != nil {
		return "", err
	}
	return val.(string), nil
}

func ZPopMax(ctx context.Context, rc *redis.Redis, key string) (string, error) {
	if ZPopMaxSha1 == "" {
		var err error
		ZPopMaxSha1, err = rc.ScriptLoadCtx(ctx, ZPopMaxScript)
		if err != nil {
			return "", err
		}
	}
	val, err := rc.EvalShaCtx(ctx, ZPopMaxSha1, []string{key})
	if err != nil {
		return "", err
	}
	return val.(string), nil
}
