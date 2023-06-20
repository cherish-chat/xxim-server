package messagemodel

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xcache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"time"
)

type xRedisSeq struct {
	rc *redis.Redis
}

var RedisSeq *xRedisSeq

func InitRedisSeq(rc *redis.Redis) {
	RedisSeq = &xRedisSeq{rc: rc}
}

type ConvSeq struct {
	ConvId     string
	ConvType   pb.ConversationType
	MinSeq     int64
	MaxSeq     int64
	UpdateTime int64
}

type ConvInfo struct {
	ConvId   string
	ConvType pb.ConversationType
}

// incrConvMaxSeqScript: hincrby key hkey 1; hset key UpdateTime now
const incrConvMaxSeqScript = `
local seq = redis.call('hincrby', KEYS[1], KEYS[2], 1)
redis.call('hset', KEYS[1], "updateTime", ARGV[1])
return seq
`

var incrConvMaxSeqSha string

func (x *xRedisSeq) IncrConvMaxSeq(ctx context.Context, convId string, convTyp pb.ConversationType) (int, error) {
	rc := x.rc
	if incrConvMaxSeqSha == "" {
		var err error
		incrConvMaxSeqSha, err = rc.ScriptLoadCtx(ctx, incrConvMaxSeqScript)
		if err != nil {
			return 0, err
		}
	}
	result, err := rc.EvalShaCtx(ctx, incrConvMaxSeqSha, []string{xcache.RedisVal.HashKeyConvKv(convId, int32(convTyp)), xcache.RedisVal.HashKeyConvMaxSeq, "updateTime"}, time.Now().UnixMilli())
	if err != nil {
		return 0, err
	}
	return int(utils.Number.Any2Int64(result)), nil
}

// getConvMaxSeqScript: hmget key hkey1 hkey2 UpdateTime
const getConvMaxSeqScript = `
local result = redis.call('hmget', KEYS[1], KEYS[2], KEYS[3], KEYS[4])
return result
`

var getConvMaxSeqSha string

func (x *xRedisSeq) GetConvMaxSeq(ctx context.Context, userId string, convId string, convTyp pb.ConversationType) (*ConvSeq, error) {
	rc := x.rc
	if getConvMaxSeqSha == "" {
		var err error
		getConvMaxSeqSha, err = rc.ScriptLoadCtx(ctx, getConvMaxSeqScript)
		if err != nil {
			return nil, err
		}
	}
	result, err := rc.EvalShaCtx(ctx, getConvMaxSeqSha, []string{xcache.RedisVal.HashKeyConvKv(convId, int32(convTyp)), xcache.RedisVal.HashKeyConvMaxSeq, xcache.RedisVal.HashKeyConvMinSeq(userId), "updateTime"})
	if err != nil {
		return nil, err
	}
	arr := result.([]interface{})
	return &ConvSeq{
		ConvId:     convId,
		ConvType:   convTyp,
		MaxSeq:     utils.Number.Any2Int64(arr[0]),
		MinSeq:     utils.Number.Any2Int64(arr[1]),
		UpdateTime: utils.Number.Any2Int64(arr[2]),
	}, nil
}

// batchGetConvMaxSeqScript: hmget key1 arg1 arg2 UpdateTime; hmget key2 arg1 arg2 UpdateTime
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

func (x *xRedisSeq) BatchGetConvMaxSeq(ctx context.Context, userId string, convIdTypes ...*ConvInfo) (map[*ConvInfo]*ConvSeq, error) {
	m := make(map[*ConvInfo]*ConvSeq, 0)
	if len(convIdTypes) == 0 {
		return m, nil
	}
	rc := x.rc
	if batchGetConvMaxSeqSha == "" {
		var err error
		batchGetConvMaxSeqSha, err = rc.ScriptLoadCtx(ctx, batchGetConvMaxSeqScript)
		if err != nil {
			return nil, err
		}
	}
	keys := make([]string, 0, len(convIdTypes))
	for _, convId := range convIdTypes {
		keys = append(keys, xcache.RedisVal.HashKeyConvKv(convId.ConvId, int32(convId.ConvType)))
	}
	result, err := rc.EvalShaCtx(ctx, batchGetConvMaxSeqSha, keys, xcache.RedisVal.HashKeyConvMaxSeq, xcache.RedisVal.HashKeyConvMinSeq(userId), "updateTime")
	if err != nil {
		return nil, err
	}
	arr := result.([]interface{})
	for i, v := range arr {
		convIdType := convIdTypes[i]
		arr := v.([]interface{})
		m[convIdType] = &ConvSeq{
			ConvId:     convIdType.ConvId,
			ConvType:   convIdType.ConvType,
			MaxSeq:     utils.Number.Any2Int64(arr[0]),
			MinSeq:     utils.Number.Any2Int64(arr[1]),
			UpdateTime: utils.Number.Any2Int64(arr[2]),
		}
	}
	return m, nil
}
