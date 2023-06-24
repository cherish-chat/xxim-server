package subscriptionmodel

import (
	opts "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SubscriptionSubscribe struct {
	SubscriptionId string             `bson:"subscriptionId" json:"subscriptionId"`
	MemberUserId   string             `bson:"memberUserId" json:"memberUserId"`
	SubscribeTime  primitive.DateTime `bson:"subscribeTime" json:"subscribeTime"`
}

func (m *SubscriptionSubscribe) GetIndexes() []opts.IndexModel {
	return []opts.IndexModel{{
		Key:          []string{"subscriptionId", "memberUserId"},
		IndexOptions: options.Index().SetUnique(true),
	}}
}
