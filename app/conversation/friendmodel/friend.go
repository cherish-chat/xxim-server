package friendmodel

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
)

// Friend 好友 数据库模型
type Friend struct {
	// UserId 用户ID 申请加好友的用户ID
	UserId string `json:"userId" gorm:"column:userId;type:char(32);primary_key;not null"`
	// FriendId 好友ID 被申请加好友的用户ID
	FriendId string `json:"friendId" gorm:"column:friendId;type:char(32);primary_key;not null"`
	// BeFriendTime 成为好友时间 13位时间戳
	BeFriendTime int64 `json:"beFriendTime" gorm:"column:beFriendTime;type:bigint(13);not null"`
}

// TableName 表名
func (m *Friend) TableName() string {
	return "friend"
}

// xFriendModel 数据库操作实例
type xFriendModel struct {
	db *gorm.DB
	rc *redis.Redis
}

var FriendModel *xFriendModel

// InitFriendModel 初始化数据库操作实例
func InitFriendModel(db *gorm.DB, rc *redis.Redis) {
	FriendModel = &xFriendModel{
		db: db,
		rc: rc,
	}
	db.AutoMigrate(&Friend{})
}
