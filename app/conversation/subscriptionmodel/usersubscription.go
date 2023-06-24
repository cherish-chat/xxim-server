package subscriptionmodel

import (
	opts "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserSubscription struct {
	SubscriptionId string             `bson:"subscriptionId" json:"subscriptionId"`
	Subscriber     string             `bson:"subscriber" json:"subscriber"`
	SubscribeTime  primitive.DateTime `bson:"subscribeTime" json:"subscribeTime"`
	ExtraMap       bson.M             `bson:"extraMap,omitempty" json:"extraMap"`
}

func (m *UserSubscription) GetIndexes() []opts.IndexModel {
	return []opts.IndexModel{{
		Key:          []string{"subscriptionId", "subscriber"},
		IndexOptions: options.Index().SetUnique(true),
	}}
}
