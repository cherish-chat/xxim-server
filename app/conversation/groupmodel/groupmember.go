package groupmodel

import (
	"github.com/qiniu/qmgo"
	opts "github.com/qiniu/qmgo/options"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GroupMember struct {
	// GroupId 群组ID
	GroupId int `gorm:"column:groupId;type:bigint(20);not null;index:idx_groupId;primary_key;" bson:"groupId"`
	// MemberUserId 群成员ID
	MemberUserId int `gorm:"column:memberUserId;type:bigint(20);not null;index:idx_memberUserId;primary_key;" bson:"memberUserId"`
	// JoinTime 加入时间
	JoinTime int `gorm:"column:joinTime;type:int(11);not null;" bson:"joinTime"`
	// JoinSource 加入来源
	JoinSource string `gorm:"column:joinSource;type:varchar(32);not null;default:'';" bson:"joinSource"`
}

func (m *GroupMember) TableName() string {
	return "group_member"
}

// GetIndexes 索引
func (m *GroupMember) GetIndexes() []opts.IndexModel {
	return []opts.IndexModel{{
		Key:          []string{"groupId", "memberUserId"},
		IndexOptions: options.Index().SetName("unique_group_member").SetUnique(true),
	}, {
		Key: []string{"groupId"},
	}, {
		Key: []string{"memberUserId"},
	}, {
		Key: []string{"joinTime"},
	}}
}

type xGroupMemberModel struct {
	coll *qmgo.QmgoClient
	rc   *redis.Redis
}

var GroupMemberModel *xGroupMemberModel

func InitGroupMemberModel(coll *qmgo.QmgoClient, rc *redis.Redis) {
	GroupMemberModel = &xGroupMemberModel{
		coll: coll,
		rc:   rc,
	}
}
