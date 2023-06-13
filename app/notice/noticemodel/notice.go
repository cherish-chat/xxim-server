package noticemodel

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

type ConversationType int8
type NoticeType int8
type ContentType int8

const (
	// ConversationTypeSingle 单聊
	ConversationTypeSingle ConversationType = iota
	// ConversationTypeGroup 群聊
	ConversationTypeGroup
	// ConversationTypeMarketing 营销号
	ConversationTypeMarketing
)

type NoticeOption struct {
	// PersistForServer 服务端是否持久化
	PersistForServer bool `bson:"persistForServer" json:"persistForServer"`
	// PersistForClient 客户端是否持久化
	PersistForClient bool `bson:"persistForClient" json:"persistForClient"`
}

// Notice 通知 数据库模型
// Notice 是一个特殊的消息，用于通知用户，例如：邀请进群了、被踢出群了、群解散了、营销号推送了一条消息等等
// 这种消息的特点 具有时效性，例如：发送了一条邀请进群的通知，但此时群已经解散了，那么这条通知就没有意义了，所以收到这种通知后，客户端应该同步服务端信息，再做判断
type Notice struct {
	//NoticeId 通知id
	NoticeId string `bson:"_id" json:"noticeId"`
	// ConversationId 会话ID
	ConversationId string `bson:"conversationId" json:"conversationId"`
	// ConversationType 会话类型
	ConversationType ConversationType `bson:"conversationType" json:"conversationType"`
	// NoticeType 通知类型
	NoticeType NoticeType `bson:"noticeType" json:"noticeType"`
	// UniqueKey 唯一标识
	UniqueKey string `bson:"uniqueKey" json:"uniqueKey"`
	// UpdateTime 更新时间 客户端将使用UpdateTime>?的方式拉取最新的通知
	UpdateTime primitive.DateTime `bson:"updateTime" json:"updateTime"`

	// Options 通知选项
	Options NoticeOption `gorm:"column:options;type:json;" json:"options"`

	// Content 通知内容
	Content []byte `bson:"content" json:"content"`
	// ContentType 通知内容类型
	ContentType ContentType `bson:"contentType" json:"contentType"`
	// Title 通知标题 显示在会话列表 进行预览
	Title string `bson:"title" json:"title"`

	// Extra 扩展字段
	Extra bson.M `bson:"extra" json:"extra"`
}

func GenerateNoticeId(conversationId string, conversationType ConversationType, noticeType NoticeType, uniqueKey string) string {
	param := url.Values{}
	param.Add("conversationId", conversationId)
	param.Add("conversationType", fmt.Sprintf("%d", conversationType))
	param.Add("noticeType", fmt.Sprintf("%d", noticeType))
	param.Add("uniqueKey", uniqueKey)
	s := param.Encode()
	return utils.Md5(s)
}

func (m *Notice) GetIndexes() []opts.IndexModel {
	return []opts.IndexModel{{
		Key:          []string{"-updateTime"},
		IndexOptions: options.Index().SetName("updateTime"),
	}, {
		Key:          []string{"-conversationId", "conversationType", "noticeType"},
		IndexOptions: nil,
	}}
}

type xNoticeModel struct {
	coll *qmgo.QmgoClient
	rc   *redis.Redis
}

var NoticeModel *xNoticeModel

func InitNoticeModel(coll *qmgo.QmgoClient, rc *redis.Redis) {
	NoticeModel = &xNoticeModel{
		coll: coll,
		rc:   rc,
	}
}
