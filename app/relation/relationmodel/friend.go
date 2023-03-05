package relationmodel

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xredis"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
)

type Friend struct {
	// 发起好友请求的用户
	UserId string `json:"userId" bson:"userId" gorm:"column:userId;type:char(32);not null;index:idx_user_friend_id,unique;index;comment:发起好友请求的用户"`
	// 被添加的用户
	FriendId string `json:"friendId" bson:"friendId" gorm:"column:friendId;type:char(32);not null;index:idx_user_friend_id,unique;index;comment:被添加的用户"`
	// 创建时间
	CreateTime int64 `json:"createTime" bson:"createTime" gorm:"column:createTime;type:bigint(20);not null;index;comment:创建时间"`
}

func (m *Friend) TableName() string {
	return "friend"
}

func GetMyFriendList(ctx context.Context, rc *redis.Redis, tx *gorm.DB, userId string) ([]string, error) {
	// 从 redis 中获取
	friends, err := getMyFriendListFromRedis(ctx, rc, userId)
	if err != nil {
		return getMyFriendListFromMysql(ctx, rc, tx, userId)
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

func getMyFriendListFromMysql(ctx context.Context, rc *redis.Redis, tx *gorm.DB, userId string) ([]string, error) {
	var friends []*Friend
	err := tx.Where("userId = ?", userId).Find(&friends).Error
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

func AreMyFriend(ctx context.Context, rc *redis.Redis, tx *gorm.DB, userId string, friendIds []string) (map[string]bool, error) {
	existMap, err := xredis.HMExist(rc, ctx, rediskey.FriendList(userId), friendIds...)
	if err != nil {
		listFromMysql, err := getMyFriendListFromMysql(ctx, rc, tx, userId)
		if err != nil {
			return nil, err
		}
		m := make(map[string]bool, len(friendIds))
		for _, friendId := range friendIds {
			m[friendId] = utils.InSlice(listFromMysql, friendId)
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
