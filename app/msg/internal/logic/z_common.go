package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	zedis "github.com/zeromicro/go-zero/core/stores/redis"
	"time"
)

type convSeq struct {
	convId     string
	minSeq     int64
	maxSeq     int64
	updateTime int64
}

// incrConvMaxSeqScript: hincrby key hkey 1; hset key updateTime now
const incrConvMaxSeqScript = `
local seq = redis.call('hincrby', KEYS[1], KEYS[2], 1)
redis.call('hset', KEYS[1], "updateTime", ARGV[1])
return seq
`

var incrConvMaxSeqSha string

func IncrConvMaxSeq(rc *zedis.Redis, ctx context.Context, convId string) (int, error) {
	if incrConvMaxSeqSha == "" {
		var err error
		incrConvMaxSeqSha, err = rc.ScriptLoadCtx(ctx, incrConvMaxSeqScript)
		if err != nil {
			return 0, err
		}
	}
	result, err := rc.EvalShaCtx(ctx, incrConvMaxSeqSha, []string{rediskey.ConvKv(convId), rediskey.HKConvMaxSeq(), "updateTime"}, time.Now().UnixMilli())
	if err != nil {
		return 0, err
	}
	return int(utils.AnyToInt64(result)), nil
}

// getConvMaxSeqScript: hmget key hkey1 hkey2 updateTime
const getConvMaxSeqScript = `
local result = redis.call('hmget', KEYS[1], KEYS[2], KEYS[3], KEYS[4])
return result
`

var getConvMaxSeqSha string

func GetConvMaxSeq(rc *zedis.Redis, ctx context.Context, userId string, convId string) (*convSeq, error) {
	if getConvMaxSeqSha == "" {
		var err error
		getConvMaxSeqSha, err = rc.ScriptLoadCtx(ctx, getConvMaxSeqScript)
		if err != nil {
			return nil, err
		}
	}
	result, err := rc.EvalShaCtx(ctx, getConvMaxSeqSha, []string{rediskey.ConvKv(convId), rediskey.HKConvMaxSeq(), rediskey.HKConvMinSeq(userId), "updateTime"})
	if err != nil {
		return nil, err
	}
	arr := result.([]interface{})
	return &convSeq{
		convId:     convId,
		maxSeq:     utils.AnyToInt64(arr[0]),
		minSeq:     utils.AnyToInt64(arr[1]),
		updateTime: utils.AnyToInt64(arr[2]),
	}, nil
}

// batchGetConvMaxSeqScript: hmget key1 arg1 arg2 updateTime; hmget key2 arg1 arg2 updateTime
const batchGetConvMaxSeqScript = `
local result = {}
for i=1, #KEYS do
	local key = KEYS[i]
	local arr = redis.call('hmget', key, ARGV[1], ARGV[2], ARGV[3])
	table.insert(result, arr)
end
return result
`

var batchGetConvMaxSeqSha string

func BatchGetConvMaxSeq(rc *zedis.Redis, ctx context.Context, userId string, convIds []string) (map[string]*convSeq, error) {
	if batchGetConvMaxSeqSha == "" {
		var err error
		batchGetConvMaxSeqSha, err = rc.ScriptLoadCtx(ctx, batchGetConvMaxSeqScript)
		if err != nil {
			return nil, err
		}
	}
	keys := make([]string, 0, len(convIds))
	for _, convId := range convIds {
		keys = append(keys, rediskey.ConvKv(convId))
	}
	result, err := rc.EvalShaCtx(ctx, batchGetConvMaxSeqSha, keys, rediskey.HKConvMaxSeq(), rediskey.HKConvMinSeq(userId), "updateTime")
	if err != nil {
		return nil, err
	}
	arr := result.([]interface{})
	m := make(map[string]*convSeq, len(arr))
	for i, v := range arr {
		convId := convIds[i]
		arr := v.([]interface{})
		m[convId] = &convSeq{
			convId:     convId,
			maxSeq:     utils.AnyToInt64(arr[0]),
			minSeq:     utils.AnyToInt64(arr[1]),
			updateTime: utils.AnyToInt64(arr[2]),
		}
	}
	return m, nil
}
