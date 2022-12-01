package xredis

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

const (
	zaddexScript = `
local key = KEYS[1]
local score = tonumber(ARGV[1])
local member = ARGV[2]
local ttl = tonumber(ARGV[3])
redis.call('zadd', key, score, member)
redis.call('expire', key, ttl)
return 1
`
)

var (
	zaddexSha1 string
)

func ZaddEx(ctx context.Context, rc *redis.Redis, key string, member string, score int64, seconds int64) error {
	if zaddexSha1 == "" {
		var err error
		mhSetSha, err = rc.ScriptLoadCtx(ctx, zaddexScript)
		if err != nil {
			return err
		}
	}
	_, err := rc.EvalShaCtx(ctx, zaddexSha1, []string{key}, score, member, seconds)
	return err
}
