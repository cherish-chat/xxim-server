package xredis

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/utils"
	zedis "github.com/zeromicro/go-zero/core/stores/redis"
)

// incrEx is a redis script for incrEx
const incrEx = `
local key = KEYS[1]
local expire = ARGV[1]
local incr = ARGV[2]
local value = redis.call("incrby", key, incr)
redis.call("expire", key, expire)
return value
`

var (
	incrExSha = ""
)

func IncrEx(rc *zedis.Redis, ctx context.Context, key string, expire int64, incr int64) (int64, error) {
	if incrExSha == "" {
		var err error
		incrExSha, err = rc.ScriptLoadCtx(ctx, incrEx)
		if err != nil {
			return 0, err
		}
	}
	result, err := rc.EvalShaCtx(ctx, incrExSha, []string{key}, expire, incr)
	if err != nil {
		return 0, err
	}
	return utils.AnyToInt64(result), nil
}
