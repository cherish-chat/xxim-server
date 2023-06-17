package conversationmodel

import (
	"github.com/qiniu/qmgo"
	opts "github.com/qiniu/qmgo/options"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ValueType int8
type ConversationType int8

const (
	// ValueTypeString 字符串
	ValueTypeString ValueType = iota
	// ValueTypeBool 布尔值
	ValueTypeBool
	// ValueTypeInt 整数
	ValueTypeInt
	// ValueTypeFloat 浮点数
	ValueTypeFloat
	// ValueTypeArrayJson 数组JSON
	ValueTypeArrayJson
	// ValueTypeMapJson MapJSON
	ValueTypeMapJson
	// ValueTypeFileUrl 文件url;
	// valueTypeExt会记录允许的拓展名;
	// example: [".png", ".jpg", ".jpeg", ".gif", ".bmp"]
	ValueTypeFileUrl
)

type ConversationSetting struct {
	// K key
	K string `bson:"k" json:"k"`
	// G group 分组
	G string `bson:"g" json:"g"`
	// V value
	V string `bson:"v" json:"v"`
	// ValueType 值类型
	ValueType ValueType `bson:"valueType" json:"valueType"`
	// ValueTypeExt 值类型扩展
	ValueTypeExt string `bson:"valueTypeExt" json:"valueTypeExt"`
}

type ConversationMember struct {
	// ConversationId 群组ID
	ConversationId int `bson:"conversationId"`
	// ConversationType 会话类型
	ConversationType ConversationType `json:"conversationType"`
	// MemberUserId 群成员ID
	MemberUserId int `bson:"memberUserId"`
	// JoinTime 加入时间
	JoinTime int `bson:"joinTime"`
	// JoinSource 加入来源
	JoinSource string `bson:"joinSource"`
	// Settings 群成员设置
	Settings []*ConversationSetting ` bson:"settings"`
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
