package xredis

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

var hmsetexSha = ""

// HMSetEx sets the hash with ttl.
func HMSetEx(rc *redis.Redis, ctx context.Context, key string, fields map[string]interface{}, ttl int64) error {
	if hmsetexSha == "" {
		var err error
		hmsetexSha, err = rc.ScriptLoadCtx(ctx, hmsetex)
		if err != nil {
			return err
		}
	}
	args := make([]interface{}, 0, len(fields)*2+2)
	args = append(args, ttl)
	for k, v := range fields {
		args = append(args, k, v)
	}
	_, err := rc.EvalShaCtx(ctx, hmsetexSha, []string{key}, args...)
	return err
}

// hmsetex lua script;
const hmsetex = `
local key = KEYS[1]
local ttl = ARGV[1]
local exists = redis.call('EXISTS', key)
redis.call('HMSET', key, unpack(ARGV, 2))
redis.call('EXPIRE', key, ttl)
return 1
` // hmsetex key ttl field value [field value ...]
