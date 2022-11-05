package logic

import (
	"context"
	"encoding/json"
	"github.com/cherish-chat/xxim-server/app/msg/msgmodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xredis"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"github.com/qiniu/qmgo"
	"github.com/zeromicro/go-zero/core/logx"
	zedis "github.com/zeromicro/go-zero/core/stores/redis"
	"go.mongodb.org/mongo-driver/bson"
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

func GetConvMaxSeq(rc *zedis.Redis, ctx context.Context, convId string) (*convSeq, error) {
	if getConvMaxSeqSha == "" {
		var err error
		getConvMaxSeqSha, err = rc.ScriptLoadCtx(ctx, getConvMaxSeqScript)
		if err != nil {
			return nil, err
		}
	}
	result, err := rc.EvalShaCtx(ctx, getConvMaxSeqSha, []string{rediskey.ConvKv(convId), rediskey.HKConvMaxSeq(), rediskey.HKConvMinSeq(), "updateTime"})
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

func BatchGetConvMaxSeq(rc *zedis.Redis, ctx context.Context, convIds []string) (map[string]*convSeq, error) {
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
	result, err := rc.EvalShaCtx(ctx, batchGetConvMaxSeqSha, keys, rediskey.HKConvMaxSeq(), rediskey.HKConvMinSeq(), "updateTime")
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

func generateSeqList(deviceMaxSeq int64, conv *convSeq) []int64 {
	// 如果 conv.MaxSeq - deviceMaxSeq > 1000, 则只返回 1000 个 seq
	if conv.maxSeq-deviceMaxSeq > 1000 {
		deviceMaxSeq = conv.maxSeq - 1000
	}
	seqList := make([]int64, 0, deviceMaxSeq-conv.minSeq+1)
	for i := conv.minSeq; i <= deviceMaxSeq; i++ {
		seqList = append(seqList, i)
	}
	return seqList
}

func MsgFromMongo(
	ctx context.Context,
	rc *zedis.Redis,
	collection *qmgo.Collection,
	ids []string,
) (msgList []*msgmodel.Msg, err error) {
	if len(ids) == 0 {
		return make([]*msgmodel.Msg, 0), nil
	}
	xtrace.StartFuncSpan(ctx, "FindMsgByIds", func(ctx context.Context) {
		err = collection.Find(ctx, bson.M{
			"_id": bson.M{"$in": ids},
		}).All(&msgList)
	})
	if err != nil {
		logx.WithContext(ctx).Errorf("GetSingleMsgListBySeq failed, err: %v", err)
		return nil, err
	}
	msgMap := make(map[string]*msgmodel.Msg)
	for _, msg := range msgList {
		msgMap[msg.ServerMsgId] = msg
		// 存入redis
		redisMsg, _ := json.Marshal(msg)
		err = rc.SetexCtx(ctx, rediskey.MsgKey(msg.ServerMsgId), string(redisMsg), msg.ExpireSeconds())
		if err != nil {
			logx.WithContext(ctx).Errorf("redis Setex error: %v", err)
			continue
		}
	}
	var notFoundIds []string
	for _, id := range ids {
		if _, ok := msgMap[id]; !ok {
			notFoundIds = append(notFoundIds, id)
		}
	}
	if len(notFoundIds) > 0 {
		// 占位符写入redis
		for _, id := range notFoundIds {
			err = rc.SetexCtx(ctx, rediskey.MsgKey(id), xredis.NotFound, xredis.ExpireMinutes(5))
			if err != nil {
				logx.WithContext(ctx).Errorf("redis Setex error: %v", err)
				continue
			}
		}
	}
	return msgList, nil
}
