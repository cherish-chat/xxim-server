package channelmodel

import (
	opts "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChannelSubscribeCache struct {
	ChannelId     string             `bson:"channelId" json:"channelId"`
	MemberUserId  string             `bson:"memberUserId" json:"memberUserId"`
	SubscribeTime primitive.DateTime `bson:"subscribeTime" json:"subscribeTime"`
}

func (m *ChannelSubscribeCache) GetIndexes() []opts.IndexModel {
	return []opts.IndexModel{{
		Key:          []string{"channelId", "memberUserId"},
		IndexOptions: options.Index().SetUnique(true),
	}, {
		Key: []string{"memberUserId"},
	}, {
		Key: []string{"channelId"},
	}, {
		Key: []string{"-subscribeTime"},
	}}
}
