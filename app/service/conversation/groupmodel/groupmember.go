package groupmodel

import (
	opts "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GroupMember struct {
	GroupId      string `bson:"groupId" json:"groupId"`
	MemberUserId string `bson:"memberUserId" json:"memberUserId"`
	// JoinTime 加入时间
	JoinTime primitive.DateTime `bson:"joinTime"`
	// JoinSource 加入来源
	JoinSource bson.M `bson:"joinSource,omitempty"`
	// Settings 群成员设置
	Settings bson.M ` bson:"settings,omitempty"`
}

// GetIndexes 获取索引
func (m *GroupMember) GetIndexes() []opts.IndexModel {
	return []opts.IndexModel{{
		Key:          []string{"groupId", "memberUserId"},
		IndexOptions: options.Index().SetUnique(true),
	}, {
		Key: []string{"memberUserId"},
	}, {
		Key: []string{"groupId"},
	}, {
		Key: []string{"-joinTime"},
	}}
}
