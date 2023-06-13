package friendmodel

import (
	"github.com/qiniu/qmgo"
	opts "github.com/qiniu/qmgo/options"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Friend 好友 数据库模型
type Friend struct {
	// UserId 用户ID 申请加好友的用户ID
	UserId string `json:"userId" gorm:"column:userId;type:char(32);primary_key;not null" bson:"userId"`
	// FriendId 好友ID 被申请加好友的用户ID
	FriendId string `json:"friendId" gorm:"column:friendId;type:char(32);primary_key;not null" bson:"friendId"`
	// BeFriendTime 成为好友时间 13位时间戳
	BeFriendTime int64 `json:"beFriendTime" gorm:"column:beFriendTime;type:bigint(13);not null" bson:"beFriendTime"`
}

// TableName 表名
func (m *Friend) TableName() string {
	return "friend"
}

// GetIndexes 获取索引
func (m *Friend) GetIndexes() []opts.IndexModel {
	return []opts.IndexModel{{
		Key:          []string{"userId", "friendId"},
		IndexOptions: options.Index().SetName("unique_friend").SetUnique(true),
	}, {
		Key: []string{"userId"},
	}}
}

// xFriendModel 数据库操作实例
type xFriendModel struct {
	coll *qmgo.QmgoClient
	rc   *redis.Redis
}

var FriendModel *xFriendModel

// InitFriendModel 初始化数据库操作实例
func InitFriendModel(coll *qmgo.QmgoClient, rc *redis.Redis) {
	FriendModel = &xFriendModel{
		coll: coll,
		rc:   rc,
	}
}
