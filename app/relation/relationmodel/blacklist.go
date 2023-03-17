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

type Blacklist struct {
	// 发起黑名单请求的用户
	UserId string `json:"userId" bson:"userId" gorm:"column:userId;type:char(32);not null;index:userId;index:userId_blacklistId,unique;comment:发起黑名单请求的用户"`
	// 被添加的用户
	BlacklistId string `json:"blacklistId" bson:"blacklistId" gorm:"column:blacklistId;type:char(32);not null;index:blacklistId;index:userId_blacklistId,unique;comment:被添加的用户"`
	// 创建时间
	CreateTime int64 `json:"createTime" bson:"createTime" gorm:"column:createTime;type:bigint(20);not null;comment:创建时间"`
}

func (m *Blacklist) TableName() string {
	return "blacklist"
}

func GetMyBlacklistList(ctx context.Context, rc *redis.Redis, tx *gorm.DB, userId string) ([]string, error) {
	// 从 redis 中获取
	blacklists, err := getMyBlacklistListFromRedis(ctx, rc, userId)
	if err != nil {
		return getMyBlacklistListFromMysql(ctx, rc, tx, userId)
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

func getMyBlacklistListFromMysql(ctx context.Context, rc *redis.Redis, tx *gorm.DB, userId string) ([]string, error) {
	var blacklists []*Blacklist
	err := tx.Where("userId = ?", userId).Find(&blacklists).Error
	if err != nil {
		return nil, err
	}
	var kvs = make(map[string]string, len(blacklists))
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

func AreMyBlacklist(ctx context.Context, rc *redis.Redis, tx *gorm.DB, userId string, blacklistIds []string) (map[string]bool, error) {
	existMap, err := xredis.HMExist(rc, ctx, rediskey.BlacklistList(userId), blacklistIds...)
	if err != nil {
		listFromMysql, err := getMyBlacklistListFromMysql(ctx, rc, tx, userId)
		listFromMysqlMap := make(map[string]bool, len(listFromMysql))
		for _, friendId := range listFromMysql {
			listFromMysqlMap[friendId] = true
		}
		if err != nil {
			return nil, err
		}
		m := make(map[string]bool, len(blacklistIds))
		for _, blacklistId := range blacklistIds {
			_, m[blacklistId] = listFromMysqlMap[blacklistId]
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
