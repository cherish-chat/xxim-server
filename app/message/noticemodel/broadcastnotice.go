package noticemodel

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	opts "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ConversationType = pb.ConversationType
type NoticeContentType = pb.NoticeContentType

// BroadcastNotice 广播通知 数据库模型
// BroadcastNotice 是一个特殊的消息，用于通知用户，例如：邀请进群了、被踢出群了、群解散了、营销号推送了一条消息等等
// 这种消息的特点 具有时效性，例如：发送了一条邀请进群的通知，但此时群已经解散了，那么这条通知就没有意义了，所以收到这种通知后，客户端应该同步服务端信息，再做判断
type BroadcastNotice struct {
	NoticeId         string             `bson:"noticeId" json:"noticeId"`
	ConversationId   string             `bson:"conversationId" json:"conversationId"`
	ConversationType ConversationType   `bson:"conversationType" json:"conversationType"`
	Content          string             `bson:"content" json:"content"`
	ContentType      NoticeContentType  `bson:"contentType" json:"contentType"`
	UpdateTime       primitive.DateTime `bson:"updateTime" json:"updateTime"`
	Sort             int64              `bson:"sort" json:"sort"`
}

func (m *BroadcastNotice) GetIndexes() []opts.IndexModel {
	return []opts.IndexModel{{
		Key:          []string{"noticeId"},
		IndexOptions: options.Index().SetUnique(true),
	}, {
		Key: []string{"-sort"},
	}, {
		Key: []string{"conversationId", "conversationType"},
	}, {
		Key:          []string{"conversationId", "conversationType", "contentType"},
		IndexOptions: options.Index().SetUnique(true),
	}}
}

func (m *BroadcastNotice) ToPb() *pb.Notice {
	return &pb.Notice{
		NoticeId:         m.NoticeId,
		ConversationId:   m.ConversationId,
		ConversationType: m.ConversationType,
		Content:          m.Content,
		ContentType:      m.ContentType,
		UpdateTime:       int64(m.UpdateTime),
		Sort:             m.Sort,
	}
}

func BroadcastNoticeFromPb(in *pb.Notice) *BroadcastNotice {
	return &BroadcastNotice{
		NoticeId:         in.NoticeId,
		ConversationId:   in.ConversationId,
		ConversationType: in.ConversationType,
		Content:          in.Content,
		ContentType:      in.ContentType,
		UpdateTime:       primitive.DateTime(in.UpdateTime),
		Sort:             in.Sort,
	}
}
