package groupmodel

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/xredis"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
	"math"
)

type (
	GroupMember struct {
		// 群id
		GroupId string `bson:"groupId" json:"groupId" gorm:"column:groupId;type:char(32);not null;index:group_user,unique;comment:群id;index;"`
		// 用户id
		UserId string `bson:"userId" json:"userId" gorm:"column:userId;type:char(32);not null;index:group_user,unique;comment:用户id;index;"`
		// 加入时间
		CreateTime int64 `bson:"createTime" json:"createTime" gorm:"column:createTime;type:bigint;not null;index;comment:加入时间"`
	}
)

func (m *GroupMember) TableName() string {
	return "group_member"
}

func ListGroupsByUserIdFromMysql(ctx context.Context, tx *gorm.DB, rc *redis.Redis, userId string) ([]string, error) {
	var groupMembers []*GroupMember
	err := tx.WithContext(ctx).Model(&GroupMember{}).Where("userId = ?", userId).Find(&groupMembers).Error
	if err != nil {
		return nil, err
	}
	// 放到redis zset中
	var groupIds []string
	var members []interface{}
	for _, group := range groupMembers {
		members = append(members, group.CreateTime, group.GroupId)
		groupIds = append(groupIds, group.GroupId)
	}
	err = xredis.ZAddsEx(rc, ctx, rediskey.GroupMemberListByUserId(userId), rediskey.GroupMemberListByUserIdExpire(), members)
	if err != nil {
		logx.Errorf("zadd group member list error: %v", err)
	}
	return groupIds, nil
}

func ListGroupsByUserIdFromRedis(ctx context.Context, tx *gorm.DB, rc *redis.Redis, userId string) ([]string, error) {
	// redis key 是否存在
	exists, err := rc.ExistsCtx(ctx, rediskey.GroupMemberListByUserId(userId))
	if err != nil {
		logx.Errorf("redis key exists error: %v", err)
		return ListGroupsByUserIdFromMysql(ctx, tx, rc, userId)
	}
	if !exists {
		return ListGroupsByUserIdFromMysql(ctx, tx, rc, userId)
	}
	// 从redis中获取
	val, err := rc.ZrangebyscoreWithScoresCtx(ctx, rediskey.GroupMemberListByUserId(userId), 0, math.MaxInt64)
	if err != nil {
		logx.Errorf("zrangebyscore error: %v", err)
		return ListGroupsByUserIdFromMysql(ctx, tx, rc, userId)
	}
	var groupIds []string
	for _, v := range val {
		groupIds = append(groupIds, v.Key)
	}
	return groupIds, nil
}

func FlushGroupsByUserIdCache(ctx context.Context, rc *redis.Redis, userIds ...string) error {
	if len(userIds) == 0 {
		return nil
	}
	var keys []string
	for _, userId := range userIds {
		keys = append(keys, rediskey.GroupMemberListByUserId(userId))
	}
	_, err := rc.DelCtx(ctx, keys...)
	return err
}

func IsGroupMember(ctx context.Context, tx *gorm.DB, rc *redis.Redis, groupId, userId string) (bool, error) {
	// redis key 是否存在
	exists, err := rc.ExistsCtx(ctx, rediskey.GroupMemberListByUserId(userId))
	if err != nil {
		logx.Errorf("redis key exists error: %v", err)
		return IsGroupMemberFromMysql(ctx, tx, rc, groupId, userId)
	}
	if !exists {
		return IsGroupMemberFromMysql(ctx, tx, rc, groupId, userId)
	}
	// 从redis中获取
	val, err := rc.ZscoreCtx(ctx, rediskey.GroupMemberListByUserId(userId), groupId)
	if err != nil {
		logx.Errorf("zscore error: %v", err)
		return IsGroupMemberFromMysql(ctx, tx, rc, groupId, userId)
	}
	return val > 0, nil
}

func IsGroupMemberFromMysql(ctx context.Context, tx *gorm.DB, rc *redis.Redis, groupId string, userId string) (bool, error) {
	var count int64
	err := tx.WithContext(ctx).Model(&GroupMember{}).Where("groupId = ? and userId = ?", groupId, userId).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
