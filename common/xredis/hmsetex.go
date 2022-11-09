package xredis

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// hmsetex lua script;
const hmsetex = `
local key = KEYS[1]
local ttl = tonumber(ARGV[1])
local fields = {}
for i = 2, #ARGV, 2 do
	table.insert(fields, ARGV[i])
end
local values = {}
for i = 3, #ARGV, 2 do
	table.insert(values, ARGV[i])
end
redis.call('HMSET', key, unpack(fields), unpack(values))
redis.call('EXPIRE', key, ttl)
`

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
