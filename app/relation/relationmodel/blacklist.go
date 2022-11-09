package relationmodel

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xredis"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go.mongodb.org/mongo-driver/bson"
	opts "go.mongodb.org/mongo-driver/mongo/options"
)

type Blacklist struct {
	UserId      string `json:"userId" bson:"userId"`           // 发起黑名单请求的用户
	BlacklistId string `json:"blacklistId" bson:"blacklistId"` // 被添加的用户
	CreateTime  int64  `json:"createTime" bson:"createTime"`   // 创建时间
}

func (m *Blacklist) CollectionName() string {
	return "blacklist"
}

func (m *Blacklist) Indexes(c *qmgo.Collection) error {
	_ = c.CreateIndexes(context.Background(), []options.IndexModel{{
		Key:          []string{"userId", "blacklistId"},
		IndexOptions: opts.Index().SetUnique(true),
	}, {
		Key:          []string{"blacklistId"},
		IndexOptions: nil,
	}, {
		Key:          []string{"userId"},
		IndexOptions: nil,
	}, {
		Key: []string{"createTime"},
	}})
	return nil
}

func GetMyBlacklistList(ctx context.Context, rc *redis.Redis, c *qmgo.Collection, userId string) ([]string, error) {
	// 从 redis 中获取
	blacklists, err := getMyBlacklistListFromRedis(ctx, rc, userId)
	if err != nil {
		return getMyBlacklistListFromMongo(ctx, rc, c, userId)
	}
	return utils.SliceRemove(blacklists, xredis.NotFound), nil
}

func getMyBlacklistListFromRedis(ctx context.Context, rc *redis.Redis, userId string) ([]string, error) {
	kv, err := xredis.HGetAll(rc, ctx, rediskey.BlacklistList(userId))
	if err != nil {
		return nil, err
	}
	var blacklists []string
	for blacklistId := range kv {
		blacklists = append(blacklists, blacklistId)
	}
	return blacklists, nil
}

func getMyBlacklistListFromMongo(ctx context.Context, rc *redis.Redis, c *qmgo.Collection, userId string) ([]string, error) {
	var blacklists []*Blacklist
	err := c.Find(ctx, bson.M{"userId": userId}).All(&blacklists)
	if err != nil {
		return nil, err
	}
	var kvs = make(map[string]any, len(blacklists))
	var blacklistIds []string
	for _, blacklist := range blacklists {
		blacklistIds = append(blacklistIds, blacklist.BlacklistId)
		kvs[blacklist.BlacklistId] = ""
	}
	if len(kvs) == 0 {
		kvs[xredis.NotFound] = ""
	}
	// hmset
	err = xredis.HMSetEx(rc, ctx, rediskey.BlacklistList(userId), kvs, rediskey.BlacklistListExpire)
	if err != nil {
		logx.WithContext(ctx).Errorf("redis hmset error: %v", err)
	}
	return blacklistIds, nil
}

func AreMyBlacklist(ctx context.Context, rc *redis.Redis, c *qmgo.Collection, userId string, blacklistIds []string) (map[string]bool, error) {
	existMap, err := xredis.HMExist(rc, ctx, rediskey.BlacklistList(userId), blacklistIds...)
	if err != nil {
		listFromMongo, err := getMyBlacklistListFromMongo(ctx, rc, c, userId)
		if err != nil {
			return nil, err
		}
		m := make(map[string]bool, len(blacklistIds))
		for _, blacklistId := range blacklistIds {
			m[blacklistId] = utils.InSlice(listFromMongo, blacklistId)
		}
		return m, nil
	}
	return existMap, nil
}

func FlushBlacklistList(ctx context.Context, rc *redis.Redis, userIds ...string) error {
	var keys []string
	for _, userId := range userIds {
		keys = append(keys, rediskey.BlacklistList(userId))
	}
	if len(keys) == 0 {
		return nil
	}
	_, err := rc.DelCtx(ctx, keys...)
	return err
}
