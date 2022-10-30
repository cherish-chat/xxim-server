package xredis

import (
	"context"
	zedis "github.com/zeromicro/go-zero/core/stores/redis"
)

// mhSet lua script
// 批量 HSET redisKeys, redisHKs, redisValues, 长度必须一致
const mhSet = `

local function tableslice(tbl, first, last, step)
	local sliced = {}
	for i = first or 1, last or #tbl, step or 1 do
		sliced[#sliced+1] = tbl[i]
	end

	return sliced
end

local redisKeys = KEYS
local redisHKs = tableslice(ARGV, 1, #ARGV/2)
local redisValues = tableslice(ARGV, #ARGV/2+1, #ARGV)
local i = 1
while (i <= #redisKeys)
do
	local hk = redisHKs[i]
	local value = redisValues[i]
	redis.call("hset", redisKeys[i], hk, value)
	i = i + 1
end
return 1
`

var (
	// mhSetSha lua script sha
	mhSetSha = ""
)

type MHSetKv struct {
	Key string
	HK  string
	V   any
}

func MHSetLua(rc *zedis.Redis, ctx context.Context, kvs ...MHSetKv) error {
	if mhSetSha == "" {
		var err error
		mhSetSha, err = rc.ScriptLoadCtx(ctx, mhSet)
		if err != nil {
			return err
		}
	}
	var redisKeys []string
	var redisHKs []any
	var redisValues []any
	for _, kv := range kvs {
		redisKeys = append(redisKeys, kv.Key)
		redisHKs = append(redisHKs, kv.HK)
		redisValues = append(redisValues, kv.V)
	}
	var args []any
	args = append(args, redisHKs...)
	args = append(args, redisValues...)
	_, err := rc.EvalShaCtx(ctx, mhSetSha, redisKeys, args...)
	return err
}
