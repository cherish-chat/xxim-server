package messagemodel

import (
	"fmt"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/qiniu/qmgo"
	opts "github.com/qiniu/qmgo/options"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/url"
)

type TargetType int8
type ContentType int8

const (
	//TargetTypeUser 用户
	TargetTypeUser TargetType = iota
	//TargetTypeGroup 群组
	TargetTypeGroup
)

// MessageOptions 消息选项
type MessageOptions struct {
	// PersistForServer 服务端是否持久化
	PersistForServer bool `bson:"persistForServer" json:"persistForServer"`
	// PersistForClient 客户端是否持久化
	PersistForClient bool `bson:"persistForClient" json:"persistForClient"`
	// Remind 是否提醒 用户不在线时是否走厂商推送
	Remind bool `bson:"remind" json:"remind"`
	// Decrypt 是否需要解密 端对端加密使用aes算法 使用 ** Diffie-Hellman密钥交换算法 ** 交换32位aesKey, aesIv=截取中间16位
	Decrypt bool `bson:"decrypt" json:"decrypt"`
	// Unread 是否未读
	Unread bool `bson:"unread" json:"unread"`
}

// Message 消息 数据库模型
type Message struct {
	// MessageId 消息id 服务端唯一
	MessageId string `bson:"_id" json:"messageId"`
	// ClientMessageId 客户端消息id 客户端唯一
	ClientMessageId string `bson:"clientMessageId" json:"clientMessageId"`
	// SenderUserId 发送者用户id
	SenderUserId string `bson:"senderUserId" json:"senderUserId"`
	// TargetType 目标类型
	TargetType TargetType `bson:"targetType" json:"targetType"`
	// TargetId 目标id
	TargetId string `bson:"targetId" json:"targetId"`
	// Seq 消息序号
	Seq int64 `bson:"seq" json:"seq"`
	// ClientTime 客户端写入的时间
	ClientTime primitive.DateTime `bson:"clientTime" json:"clientTime"`
	// CreateTime 服务端写入的时间
	CreateTime primitive.DateTime `bson:"createTime" json:"createTime"`
	// UpdateTime 服务端更新的时间
	UpdateTime primitive.DateTime `bson:"updateTime" json:"updateTime"`

	// Content 消息内容
	Content []byte `bson:"content" json:"content"`
	// ContentType 消息内容类型
	ContentType ContentType `bson:"contentType" json:"contentType"`

	// SenderInfo 发送者信息
	SenderInfo bson.M `bson:"senderInfo" json:"senderInfo"`
	// RemindUserIds 提醒用户id列表 仅群组消息有效
	RemindUserIds []string `bson:"remindUserIds" json:"remindUserIds"`
	// Options 消息选项
	Options MessageOptions `bson:"options" json:"options"`
	// Extra 扩展字段
	Extra bson.M `bson:"extra" json:"extra"`
}

func (m *Message) Indexes() []opts.IndexModel {
	return []opts.IndexModel{{
		Key:          []string{"clientMessageId"},
		IndexOptions: options.Index().SetName("clientMessageId"),
	}}
}

func GenerateMessageId(senderUserId string, targetId string, targetType TargetType, seq int64) string {
	param := url.Values{}
	param.Add("senderUserId", senderUserId)
	param.Add("targetId", targetId)
	param.Add("targetType", fmt.Sprintf("%d", targetType))
	param.Add("seq", fmt.Sprintf("%d", seq))
	s := param.Encode()
	return utils.Md5(s)
}

type xMessageModel struct {
	coll *qmgo.QmgoClient
	rc   *redis.Redis
}

var MessageModel *xMessageModel

func InitMessageModel(coll *qmgo.QmgoClient, rc *redis.Redis) {
	MessageModel = &xMessageModel{
		coll: coll,
		rc:   rc,
	}
}
