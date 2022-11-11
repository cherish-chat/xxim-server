package xredis

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// hmexist lua script; 先判断key是否存在 不存在返回-1 存在则返回存在的字段
const hmexist = `
local key = KEYS[1]
local fields = ARGV
local exists = redis.call('EXISTS', key)
if exists == 1 then
    local existFields = {}
	for i, field in ipairs(fields) do
		local exist = redis.call('HEXISTS', key, field)
		if exist == 1 then
			table.insert(existFields, field)
		end
	end
	return existFields
else
	return -1
end
`

var hmexistSha = ""

// HMExist returns the exist fields of the hash.
func HMExist(rc *redis.Redis, ctx context.Context, key string, fields ...string) (map[string]bool, error) {
	if hmexistSha == "" {
		var err error
		hmexistSha, err = rc.ScriptLoadCtx(ctx, hmexist)
		if err != nil {
			return nil, err
		}
	}
	val, err := rc.EvalShaCtx(ctx, hmexistSha, []string{key}, fields)
	if err != nil {
		return nil, err
	}
	if val == int64(-1) {
		return nil, redis.Nil
	}
	existFields, ok := val.([]interface{})
	if !ok {
		return nil, redis.Nil
	}
	m := make(map[string]bool, len(existFields))
	for _, field := range existFields {
		f, ok := field.(string)
		if !ok {
			return nil, redis.Nil
		}
		m[f] = true
	}
	for _, field := range fields {
		if _, ok := m[field]; !ok {
			m[field] = false
		}
	}
	return m, nil
}
