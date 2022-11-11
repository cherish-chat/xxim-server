package xredis

import (
	"context"
	zedis "github.com/zeromicro/go-zero/core/stores/redis"
	"strconv"
)

// MZAddEx script: ZAdd arg1 key1 key2 && EXPIRE arg1 key3; ZAdd arg2 key1 key2 && EXPIRE arg2 key3;...
const mzAddExScript = `
local score = tonumber(KEYS[1])
local member = KEYS[2]
local ttl = tonumber(KEYS[3])
for i=1, #ARGV do
	local key = ARGV[i]
	redis.call('ZADD', key, score, member)
	redis.call('EXPIRE', key, ttl)
end
return 1
`

var mzAddExSha string

func MZAddEx(rc *zedis.Redis, ctx context.Context, keys []string, score int64, member string, expireSeconds int) error {
	if mzAddExSha == "" {
		var err error
		mzAddExSha, err = rc.ScriptLoadCtx(ctx, mzAddExScript)
		if err != nil {
			return err
		}
	}
	args := make([]any, 0)
	for _, key := range keys {
		args = append(args, key)
	}
	_, err := rc.EvalShaCtx(ctx, mzAddExSha, []string{strconv.FormatInt(score, 10), member, strconv.Itoa(expireSeconds)}, args...)
	return err
}

// MZRem script: ZRem arg1 key1; ZRem arg2 key1;...
const mzRemScript = `
for i=1, #ARGV do
	local key = ARGV[i]
	redis.call('ZREM', key, KEYS[1])
end
return 1
`

var mzRemSha string

func MZRem(rc *zedis.Redis, ctx context.Context, keys []string, member string) error {
	if mzRemSha == "" {
		var err error
		mzRemSha, err = rc.ScriptLoadCtx(ctx, mzRemScript)
		if err != nil {
			return err
		}
	}
	args := make([]any, 0)
	for _, key := range keys {
		args = append(args, key)
	}
	_, err := rc.EvalShaCtx(ctx, mzRemSha, []string{member}, args...)
	return err
}
