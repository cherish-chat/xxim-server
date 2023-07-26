package relationmodel

import (
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
)

type UserRemark struct {
	UserId   string `gorm:"column:userId;type:char(32);primary_key;comment:用户ID" json:"userId"`
	TargetId string `gorm:"column:targetId;type:char(32);primary_key;comment:目标用户ID" json:"targetId"`
	Remark   string `gorm:"column:remark;type:varchar(64);comment:备注" json:"remark"`
}

func (m *UserRemark) TableName() string {
	return "user_remark"
}

func FlushUserRemarkCache(rc *redis.Redis, userId string) {
	_, _ = rc.Del(rediskey.UserRemarkHashKey(userId))
}

func GetUserRemarkMap(rc *redis.Redis, tx *gorm.DB, userId string) (map[string]string, error) {
	m := make(map[string]string)
	// exists
	exists, err := rc.Exists(rediskey.UserRemarkHashKey(userId))
	if err != nil {
		return GetUserRemarkMapFromMysql(rc, tx, userId)
	}
	if !exists {
		return GetUserRemarkMapFromMysql(rc, tx, userId)
	}
	// hgetall
	allMap, err := rc.Hgetall(rediskey.UserRemarkHashKey(userId))
	if err != nil {
		return m, err
	}
	return allMap, nil
}

func GetUserRemarkMapFromMysql(rc *redis.Redis, tx *gorm.DB, userId string) (map[string]string, error) {
	m := make(map[string]string)
	var list []UserRemark
	if err := tx.Model(&UserRemark{}).Where("userId = ?", userId).Find(&list).Error; err != nil {
		return m, err
	}
	for _, v := range list {
		m[v.TargetId] = v.Remark
	}
	// hmset
	if len(m) == 0 {
		return m, nil
	}
	if err := rc.Hmset(rediskey.UserRemarkHashKey(userId), m); err != nil {
		return m, err
	}
	return m, nil
}
