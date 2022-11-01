package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	zedis "github.com/zeromicro/go-zero/core/stores/redis"
	"time"
)

// incrConvMaxSeqScript: hincrby key hkey 1; hset key updateTime now
const incrConvMaxSeqScript = `
local seq = redis.call('hincrby', KEYS[1], KEYS[2], 1)
redis.call('hset', KEYS[1], "updateTime", ARGV[1])
return seq
`

var incrConvMaxSeqSha1 string

func IncrConvMaxSeq(rc *zedis.Redis, ctx context.Context, convId string) (int, error) {
	if incrConvMaxSeqSha1 == "" {
		var err error
		incrConvMaxSeqSha1, err = rc.ScriptLoadCtx(ctx, incrConvMaxSeqScript)
		if err != nil {
			return 0, err
		}
	}
	result, err := rc.EvalShaCtx(ctx, incrConvMaxSeqSha1, []string{rediskey.ConvKv(convId), rediskey.HKConvMaxSeq(), "updateTime"}, time.Now().UnixMilli())
	if err != nil {
		return 0, err
	}
	return int(utils.AnyToInt64(result)), nil
}
