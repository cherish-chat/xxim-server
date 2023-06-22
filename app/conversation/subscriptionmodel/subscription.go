package subscriptionmodel

import (
	opts "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	ConversationIdFriendNotification = "friend_notification" // 好友通知
)

type SubscriptionType int8

const (
	SubscriptionTypeSystem      SubscriptionType = iota // 系统通知, 建议显示在会话列表, 接收系统推送来的消息, 比如停机维护...
	SubscriptionTypeInteractive                         // 互动消息, 建议显示在会话列表, 接收用户互动消息, 比如点赞、评论、回复...
	SubscriptionTypePublic                              // 公众号, 建议显示在会话列表, 接收公众号推送来的消息, 比如营销号推送的消息...
	SubscriptionTypeHidden                              // 隐藏号, 不显示在会话列表, 接收隐藏号推送来的消息, 比如上线通知、下线通知、朋友圈更新通知、昵称头像更新通知...
)

// Subscription 订阅号 数据库模型
// Subscription与Friend、Group同级，都是会话，但是Subscription是单向的，只有订阅号可以给用户发送消息，用户不能给订阅号发送消息
// 场景1：系统创建的订阅号，比如：系统通知、互动消息、活动消息等
// 场景2：用户为了向订阅者推送通知（上线通知、下线通知、朋友圈更新通知、昵称头像更新通知等）而创建的订阅号
type Subscription struct {
	SubscriptionType SubscriptionType `bson:"subscriptionType" json:"subscriptionType"` // 订阅号类型
	SubscriptionId   string           `bson:"subscriptionId" json:"subscriptionId"`     // 订阅号ID
	// 头像
	Avatar string `bson:"avatar" json:"avatar"`
	// 昵称
	Nickname string `bson:"nickname" json:"nickname"`
	// 简介
	Bio string `bson:"bio" json:"bio"`
	// 封面
	Cover string `bson:"cover" json:"cover"`
	// 扩展字段
	ExtraMap bson.M `bson:"extraMap,omitempty" json:"extraMap"`
}

func (m *Subscription) GetIndexes() []opts.IndexModel {
	return []opts.IndexModel{
		{
			Key:          []string{"subscriptionId", "subscriptionType"},
			IndexOptions: options.Index().SetUnique(true),
		},
	}
}
