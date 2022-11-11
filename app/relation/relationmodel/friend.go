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

type Friend struct {
	UserId   string `json:"userId" bson:"userId"`     // 发起好友请求的用户
	FriendId string `json:"friendId" bson:"friendId"` // 被添加的用户
}

func (m *Friend) CollectionName() string {
	return "friend"
}

func (m *Friend) Indexes(c *qmgo.Collection) error {
	_ = c.CreateIndexes(context.Background(), []options.IndexModel{{
		Key:          []string{"userId", "friendId"},
		IndexOptions: opts.Index().SetUnique(true),
	}, {
		Key:          []string{"friendId"},
		IndexOptions: nil,
	}, {
		Key:          []string{"userId"},
		IndexOptions: nil,
	}})
	return nil
}

func GetMyFriendList(ctx context.Context, rc *redis.Redis, c *qmgo.Collection, userId string) ([]string, error) {
	// 从 redis 中获取
	friends, err := getMyFriendListFromRedis(ctx, rc, userId)
	if err != nil {
		return getMyFriendListFromMongo(ctx, rc, c, userId)
	}
	return utils.SliceRemove(friends, xredis.NotFound), nil
}

func getMyFriendListFromRedis(ctx context.Context, rc *redis.Redis, userId string) ([]string, error) {
	kv, err := xredis.HGetAll(rc, ctx, rediskey.FriendList(userId))
	if err != nil {
		return nil, err
	}
	var friends []string
	for friendId := range kv {
		friends = append(friends, friendId)
	}
	return friends, nil
}

func getMyFriendListFromMongo(ctx context.Context, rc *redis.Redis, c *qmgo.Collection, userId string) ([]string, error) {
	var friends []*Friend
	err := c.Find(ctx, bson.M{"userId": userId}).All(&friends)
	if err != nil {
		return nil, err
	}
	var kvs = make(map[string]any, len(friends))
	var friendIds []string
	for _, friend := range friends {
		friendIds = append(friendIds, friend.FriendId)
		kvs[friend.FriendId] = ""
	}
	if len(kvs) == 0 {
		kvs[xredis.NotFound] = ""
	}
	// hmset
	err = xredis.HMSetEx(rc, ctx, rediskey.FriendList(userId), kvs, rediskey.FriendListExpire)
	if err != nil {
		logx.WithContext(ctx).Errorf("redis hmset error: %v", err)
	}
	return friendIds, nil
}

func AreMyFriend(ctx context.Context, rc *redis.Redis, c *qmgo.Collection, userId string, friendIds []string) (map[string]bool, error) {
	existMap, err := xredis.HMExist(rc, ctx, rediskey.FriendList(userId), friendIds...)
	if err != nil {
		listFromMongo, err := getMyFriendListFromMongo(ctx, rc, c, userId)
		if err != nil {
			return nil, err
		}
		m := make(map[string]bool, len(friendIds))
		for _, friendId := range friendIds {
			m[friendId] = utils.InSlice(listFromMongo, friendId)
		}
		return m, nil
	}
	return existMap, nil
}

func FlushFriendList(ctx context.Context, rc *redis.Redis, userIds ...string) error {
	var keys []string
	for _, userId := range userIds {
		keys = append(keys, rediskey.FriendList(userId))
	}
	if len(keys) == 0 {
		return nil
	}
	_, err := rc.DelCtx(ctx, keys...)
	return err
}
