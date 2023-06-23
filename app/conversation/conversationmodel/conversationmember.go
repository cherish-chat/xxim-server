package conversationmodel

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/qiniu/qmgo"
	opts "github.com/qiniu/qmgo/options"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ConversationType = pb.ConversationType
type ConversationSettingKey = pb.ConversationSettingKey

type ConversationMember struct {
	// ConversationId 群组ID
	ConversationId string `bson:"conversationId"`
	// ConversationType 会话类型
	ConversationType ConversationType `json:"conversationType"`
	// MemberUserId 群成员ID
	MemberUserId string `bson:"memberUserId"`
	// JoinTime 加入时间
	JoinTime primitive.DateTime `bson:"joinTime"`
	// JoinSource 加入来源
	JoinSource bson.M `bson:"joinSource,omitempty"`
	// Settings 群成员设置
	Settings bson.M ` bson:"settings,omitempty"`
}

func (m *ConversationMember) TableName() string {
	return "conversation_member"
}

// GetIndexes 索引
func (m *ConversationMember) GetIndexes() []opts.IndexModel {
	return []opts.IndexModel{{
		Key:          []string{"conversationId", "conversationType", "memberUserId"},
		IndexOptions: options.Index().SetName("unique_conversation_member").SetUnique(true),
	}, {
		Key: []string{"conversationId"},
	}, {
		Key: []string{"memberUserId"},
	}, {
		Key: []string{"joinTime"},
	}}
}

type xConversationMemberModel struct {
	coll *qmgo.QmgoClient
	rc   *redis.Redis
}

var ConversationMemberModel *xConversationMemberModel

func InitConversationMemberModel(coll *qmgo.QmgoClient, rc *redis.Redis) {
	ConversationMemberModel = &xConversationMemberModel{
		coll: coll,
		rc:   rc,
	}
}
