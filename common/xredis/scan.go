package xredis

import (
	"context"
	"github.com/go-redis/redis/v8"
)

const (
	scriptHMSetEx = `
local rediskey = KEYS[1]
local expire = ARGV[1]
local args = {}
local i = 2
while (i <= #ARGV)
do
	table.insert(args, ARGV[i])
	table.insert(args, ARGV[i+1])
	i = i + 2
end
redis.call("hmset", rediskey, unpack(args))
redis.call("expire", rediskey, expire)
return 1
`
)

var (
	scriptHMSetExSha = ""
)

func HMSetEx(rc redis.UniversalClient, ctx context.Context, rediskey string, kv map[string]any, ex int64) {
	if scriptHMSetExSha == "" {
		var err error
		scriptHMSetExSha, err = redis.NewScript(scriptHMSetEx).Load(ctx, rc).Result()
		if err != nil {
			panic(err)
		}
	}
	var args []interface{}
	args = append(args, ex)
	for k, v := range kv {
		args = append(args, k)
		args = append(args, v)
	}
	_, err := rc.EvalSha(ctx, scriptHMSetExSha, []string{rediskey}, args...).Result()
	if err != nil {
		panic(err)
	}
}

func Scan(rc redis.UniversalClient, ctx context.Context, key string, keys chan []string) {
	defer func() {
		close(keys)
	}()
	cursor := uint64(0)
	for {
		var res []string
		var e error
		res, cursor, e = rc.Scan(ctx, cursor, key, 10000).Result()
		keys <- res
		if e != nil {
			return
		}
		if cursor == 0 {
			keys <- []string{"@@end@@"}
			break
		}
	}
}
