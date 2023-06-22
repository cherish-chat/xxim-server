package noticemodel

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	opts "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SubscriptionNotice 订阅号通知 数据库模型
type SubscriptionNotice struct {
	UserId         string             `bson:"userId" json:"userId"`
	SubscriptionId string             `bson:"subscriptionId" json:"subscriptionId"`
	ContentType    NoticeContentType  `bson:"contentType" json:"contentType"`
	ContentId      string             `bson:"contentId" json:"contentId"`
	UpdateTime     primitive.DateTime `bson:"updateTime" json:"updateTime"`
	Sort           int64              `bson:"sort" json:"sort"`
}

type SubscriptionNoticeContent struct {
	ContentId string `bson:"contentId" json:"contentId"`
	Content   string `bson:"content" json:"content"`
}

func (m *SubscriptionNotice) GetIndexes() []opts.IndexModel {
	return []opts.IndexModel{{
		Key:          []string{"subscriptionId", "userId", "contentType"},
		IndexOptions: options.Index().SetUnique(true),
	}, {
		Key: []string{"-sort"},
	}, {
		Key: []string{"subscriptionId"},
	}, {
		Key:          []string{"noticeType"},
		IndexOptions: nil,
	}}
}

func (m *SubscriptionNotice) ToPb(content string) *pb.Notice {
	return &pb.Notice{
		//NoticeId:         "",
		ConversationId:   m.SubscriptionId,
		ConversationType: pb.ConversationType_Subscription,
		Content:          content,
		ContentType:      m.ContentType,
		UpdateTime:       int64(m.UpdateTime),
		Sort:             m.Sort,
	}
}
