package xredis

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// hlen lua script; 先判断key是否存在 存在则返回字段长度 否则返回-1
const hlen = `
local exists = redis.call('EXISTS', KEYS[1])
if exists == 1 then
	return redis.call('HLEN', KEYS[1])
else
	return -1
end
`

var hlenSha1 string

func HLen(rc *redis.Redis, ctx context.Context, key string) (int64, error) {
	if hlenSha1 == "" {
		var err error
		hlenSha1, err = rc.ScriptLoadCtx(ctx, hlen)
		if err != nil {
			return 0, err
		}
	}
	val, err := rc.EvalShaCtx(ctx, hlenSha1, []string{key})
	if err != nil {
		return 0, err
	}
	if val == int64(-1) {
		return 0, redis.Nil
	}
	return val.(int64), nil
}
