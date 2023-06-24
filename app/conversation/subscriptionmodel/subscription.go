package subscriptionmodel

import (
	"context"
	"github.com/qiniu/qmgo"
	opts "github.com/qiniu/qmgo/options"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

const (
	UserDefaultSubscriptionIdPrefix = "uds_"
)

const (
	//ConversationIdFriendHelper 好友助手
	ConversationIdFriendHelper = "friend_helper"
	//ConversationIdGroupHelper 群助手
	ConversationIdGroupHelper = "group_helper"
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
			Key:          []string{"subscriptionId"},
			IndexOptions: options.Index().SetUnique(true),
		},
	}
}

func InitSystemSubscription(coll *qmgo.QmgoClient) {
	subs := []*Subscription{
		{
			SubscriptionType: SubscriptionTypeSystem,
			SubscriptionId:   ConversationIdFriendHelper,
			Avatar:           "/images/system/avatar/friend_helper.png",
			Nickname:         "好友助手",
			Bio:              "该账号为系统账号，用于向您推送有关好友的系统通知",
			Cover:            "/images/system/cover/friend_helper.png",
			ExtraMap: map[string]interface{}{
				"friendHelper": "true",
			},
		},
		{
			SubscriptionType: SubscriptionTypeSystem,
			SubscriptionId:   ConversationIdGroupHelper,
			Avatar:           "/images/system/avatar/group_helper.png",
			Nickname:         "群助手",
			Bio:              "群助手",
			Cover:            "/images/system/cover/group_helper.png",
			ExtraMap: map[string]interface{}{
				"groupHelper": "true",
			},
		},
	}
	for _, sub := range subs {
		//如果查询不到则插入
		count, err := coll.Find(context.Background(), bson.M{
			"subscriptionId":   sub.SubscriptionId,
			"subscriptionType": sub.SubscriptionType,
		}).Count()
		if err != nil {
			logx.Errorf("InitSystemSubscription Find error: %v", err)
			os.Exit(1)
		}
		if count == 0 {
			logx.Infof("InitSystemSubscription Insert: %v", sub)
			_, err = coll.Upsert(context.Background(), bson.M{
				"subscriptionId":   sub.SubscriptionId,
				"subscriptionType": sub.SubscriptionType,
			}, sub, opts.UpsertOptions{
				ReplaceOptions: options.Replace().SetUpsert(true),
			})
			if err != nil {
				logx.Errorf("InitSystemSubscription Upsert error: %v", err)
				os.Exit(1)
			}
		} else {
			logx.Infof("InitSystemSubscription Find: %v", sub)
		}
	}
}

func UserSubscribedSystemConversationIds() []string {
	return []string{
		ConversationIdFriendHelper,
		ConversationIdGroupHelper,
	}
}

func UserDefaultSubscriptionId(uid string) string {
	return UserDefaultSubscriptionIdPrefix + uid
}
