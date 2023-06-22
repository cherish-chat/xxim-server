package usermodel

import (
	opts "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ValueType int32

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

// UserSetting 用户设置 数据库模型
type UserSetting struct {
	UserId string    `bson:"userId" json:"userId"`
	K      string    `bson:"k" json:"k"`
	V      string    `bson:"v" json:"v"`
	VT     ValueType `bson:"vt" json:"vt"`
}

func (m *UserSetting) GetIndexes() []opts.IndexModel {
	return []opts.IndexModel{{
		Key:          []string{"userId", "k"},
		IndexOptions: options.Index().SetUnique(true),
	}}
}

//UserSettingKey 用户设置key

const (
	//UserSettingKeyFriendApply 好友申请设置
	UserSettingKeyFriendApply = "friendApply"
	//UserSettingKeyAllowBeInvitedGroup 允许被邀请进群
	UserSettingKeyAllowBeInvitedGroup = "allowBeInvitedGroup"
)

type UserSettingFriendApplyType int32

//UserSettingFriendApplyType 用户设置好友申请类型
//好友申请验证方式
//0: 任何人都可以添加我为好友
//1: 需要发送验证消息
//2: 需要正确回答问题
//3: 需要回答问题并由我确认（不需要正确回答问题）
//4: 需要回答问题并由我确认（需要正确回答问题）
//5: 不允许任何人添加我

const (
	UserSettingFriendApplyTypeAny = iota
	UserSettingFriendApplyTypeVerifyMessage
	UserSettingFriendApplyTypeAnswerQuestion
	UserSettingFriendApplyTypeAnswerQuestionAndConfirm
	UserSettingFriendApplyTypeAnswerQuestionAndConfirmWithRightAnswer
	UserSettingFriendApplyTypeNoOne
)

// UserSettingFriendApply 用户设置好友申请
type UserSettingFriendApply struct {
	//好友申请验证方式
	//0: 任何人都可以添加我为好友
	//1: 需要发送验证消息
	//2: 需要正确回答问题
	//3: 需要回答问题并由我确认（不需要正确回答问题）
	//4: 需要回答问题并由我确认（需要正确回答问题）
	//5: 不允许任何人添加我
	ApplyType UserSettingFriendApplyType `json:"applyType"`
	//设置的问题
	Question string `json:"question"`
	//设置的答案
	Answer string `json:"answer"`
}
