package groupmodel

import (
	opts "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type GroupSubscribe struct {
	GroupId       string             `bson:"groupId" json:"groupId"`
	MemberUserId  string             `bson:"memberUserId" json:"memberUserId"`
	SubscribeTime primitive.DateTime `bson:"subscribeTime" json:"subscribeTime"`
}

func (m *GroupSubscribe) GetIndexes() []opts.IndexModel {
	return []opts.IndexModel{{
		Key:          []string{"groupId", "memberUserId"},
		IndexOptions: options.Index().SetUnique(true),
	}}
}
