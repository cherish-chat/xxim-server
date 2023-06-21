package messagemodel

import (
	"fmt"
	"github.com/cherish-chat/xxim-server/common/pb"
	opts "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ContentType = pb.MessageContentType
type ConversationType = pb.ConversationType

// MessageSender 发送者
type MessageSender struct {
	//发送者id
	Id string `bson:"id" json:"id"`
	//发送者名称
	Name string `bson:"name" json:"name"`
	//发送者头像
	Avatar string `bson:"avatar" json:"avatar"`
	//extra
	Extra string `bson:"extra" json:"extra"`
}

// MessageOptions 选项
type MessageOptions struct {
	//服务端是否保存该消息
	StorageForServer bool `bson:"storageForServer" json:"storageForServer"`
	//客户端是否保存该消息
	StorageForClient bool `bson:"storageForClient" json:"storageForClient"`
	//是否需要解密 （端对端加密技术，服务端无法解密）
	NeedDecrypt bool `bson:"needDecrypt" json:"needDecrypt"`
	//消息是否需要计入未读数
	CountUnread bool `bson:"countUnread" json:"countUnread"`
}

// Message 消息 数据库模型
type Message struct {
	// MessageId 消息id 由服务端插入时生成
	MessageId string `bson:"_id" json:"messageId"`
	// Uuid 客户端生成的id 由客户端生成 在客户端保证唯一性
	Uuid string `bson:"uuid" json:"uuid"`
	// ConversationId 发送到哪个会话
	ConversationId string `bson:"conversationId" json:"conversationId"`
	// ConversationType 会话类型
	ConversationType ConversationType `bson:"conversationType" json:"conversationType"`
	// Sender 发送者
	Sender MessageSender `bson:"sender" json:"sender"`
	// Content 消息内容
	Content []byte `bson:"content" json:"content"`
	// ContentType 消息类型
	ContentType ContentType `bson:"contentType" json:"contentType"`
	// SendTime 发送时间 由客户端生成
	SendTime primitive.DateTime `bson:"sendTime" json:"sendTime"`
	// InsertTime 插入时间 由服务端生成
	InsertTime primitive.DateTime `bson:"insertTime" json:"insertTime"`
	// Seq 在会话中的消息顺序
	Seq int64 `bson:"seq" json:"seq"`
	// Option 选项
	Option MessageOptions `bson:"option" json:"option"`
	// ExtraMap example: {"platformSource": "windows"}
	ExtraMap bson.M `bson:"extraMap" json:"extraMap"`
}

func (m *Message) GetIndexes() []opts.IndexModel {
	return []opts.IndexModel{{
		Key:          []string{"clientMessageId"},
		IndexOptions: options.Index().SetName("clientMessageId"),
	}}
}

func MessageFromPb(in *pb.Message) *Message {
	extraMap := make(bson.M)
	for k, v := range in.GetExtraMap() {
		extraMap[k] = v
	}
	return &Message{
		MessageId:        in.GetMessageId(),
		Uuid:             in.GetUuid(),
		ConversationId:   in.GetConversationId(),
		ConversationType: in.GetConversationType(),
		Sender: MessageSender{
			Id:     in.GetSender().GetId(),
			Name:   in.GetSender().GetName(),
			Avatar: in.GetSender().GetAvatar(),
			Extra:  in.GetSender().GetExtra(),
		},
		Content:     in.Content,
		ContentType: in.ContentType,
		SendTime:    primitive.DateTime(in.SendTime),
		InsertTime:  primitive.DateTime(in.InsertTime),
		Seq:         in.Seq,
		Option: MessageOptions{
			StorageForServer: in.GetOption().GetStorageForServer(),
			StorageForClient: in.GetOption().GetStorageForClient(),
			NeedDecrypt:      in.GetOption().GetNeedDecrypt(),
			CountUnread:      in.GetOption().GetCountUnread(),
		},
		ExtraMap: extraMap,
	}
}

func (m *Message) ToPb() *pb.Message {
	extraMap := make(map[string]string)
	for k, v := range m.ExtraMap {
		extraMap[k] = v.(string)
	}
	return &pb.Message{
		MessageId:        m.MessageId,
		Uuid:             m.Uuid,
		ConversationId:   m.ConversationId,
		ConversationType: m.ConversationType,
		Sender: &pb.Message_Sender{
			Id:     m.Sender.Id,
			Name:   m.Sender.Name,
			Avatar: m.Sender.Avatar,
			Extra:  m.Sender.Extra,
		},
		Content:     m.Content,
		ContentType: m.ContentType,
		SendTime:    int64(m.SendTime),
		InsertTime:  int64(m.InsertTime),
		Seq:         m.Seq,
		Option: &pb.Message_Option{
			StorageForServer: m.Option.StorageForServer,
			StorageForClient: m.Option.StorageForClient,
			NeedDecrypt:      m.Option.NeedDecrypt,
			CountUnread:      m.Option.CountUnread,
		},
		//ExtraMap: utils.Json.MarshalToString(extraMap),
		ExtraMap: extraMap,
	}
}

func (m *Message) GenerateMessageId() {
	if m.MessageId == "" {
		m.MessageId = fmt.Sprintf("%s@%d:%d", m.ConversationId, m.ConversationType, m.Seq)
	}
}
