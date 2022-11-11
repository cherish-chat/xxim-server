package xredis

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// hgetall lua script; 先判断key是否存在 存在则返回hash所有字段和值 否则返回-1
const hgetall = `
local key = KEYS[1]
local exists = redis.call('EXISTS', key)
if exists == 1 then
	return redis.call('HGETALL', key)
else
	return -1
end
`

var hgetallSha = ""

func HGetAll(rc *redis.Redis, ctx context.Context, key string) (map[string]string, error) {
	if hgetallSha == "" {
		var err error
		hgetallSha, err = rc.ScriptLoadCtx(ctx, hgetall)
		if err != nil {
			return nil, err
		}
	}
	val, err := rc.EvalShaCtx(ctx, hgetallSha, []string{key})
	if err != nil {
		return nil, err
	}
	if val == int64(-1) {
		return nil, redis.Nil
	}
	kvs, ok := val.([]interface{})
	if !ok {
		return nil, redis.Nil
	}
	if len(kvs)%2 != 0 {
		return nil, redis.Nil
	}
	m := make(map[string]string, len(kvs)/2)
	for i := 0; i < len(kvs); i += 2 {
		k, ok := kvs[i].(string)
		if !ok {
			return nil, redis.Nil
		}
		v, ok := kvs[i+1].(string)
		if !ok {
			return nil, redis.Nil
		}
		m[k] = v
	}
	return m, nil
}
