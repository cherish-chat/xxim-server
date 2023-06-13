package conversationmodel

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
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

const (
	// ConversationTypeSingle 单聊
	ConversationTypeSingle ConversationType = iota
	// ConversationTypeGroup 群聊
	ConversationTypeGroup
	// ConversationTypeMarketing 营销号
	ConversationTypeMarketing
)

// ConversationSetting 用户的会话设置 数据库模型
type ConversationSetting struct {
	// ConversationId 会话ID
	ConversationId string `gorm:"column:conversationId;type:char(32);primary_key;not null"`
	// ConversationType 会话类型
	ConversationType ConversationType `gorm:"column:conversationType;type:tinyint(1);not null;primary_key;"`
	// MemberId 成员ID
	MemberId string `gorm:"column:memberId;type:char(32);primary_key;not null"`
	// K key
	K string `gorm:"column:k;type:varchar(32);primary_key;not null"`
	// G group 分组
	G string `gorm:"column:g;type:varchar(32);not null;default:'';"`
	// V value
	V string `gorm:"column:v;type:text;"`
	// ValueType 值类型
	ValueType ValueType `gorm:"column:valueType;type:tinyint(1);not null;default:0;"`
	// ValueTypeExt 值类型扩展
	ValueTypeExt string `gorm:"column:valueTypeExt;type:varchar(32);not null;default:'';"`
}

// TableName 表名
func (m *ConversationSetting) TableName() string {
	return "conversation_setting"
}

// xConversationSettingModel 数据库操作实例
type xConversationSettingModel struct {
	db *gorm.DB
	rc *redis.Redis
}

var ConversationSettingModel *xConversationSettingModel

// InitConversationSettingModel 初始化数据库操作实例
func InitConversationSettingModel(db *gorm.DB, rc *redis.Redis) {
	ConversationSettingModel = &xConversationSettingModel{
		db: db,
		rc: rc,
	}
	db.AutoMigrate(&ConversationSetting{})
}
