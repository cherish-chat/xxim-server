package usermodel

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"gorm.io/gorm"
)

type UserSetting struct {
	UserId string            `json:"userId" bson:"userId" gorm:"column:userId;type:char(32);not null;index;index:idx_userId_key,unique"`
	Key    pb.UserSettingKey `json:"key" bson:"key" gorm:"column:key;type:varchar(64);not null;index:idx_userId_key,unique"`
	Value  string            `json:"value" bson:"value" gorm:"column:value;type:varchar(255);not null"`
}

func (m *UserSetting) TableName() string {
	return "user_setting"
}

func InitUserSetting(tx *gorm.DB) {
	tx.AutoMigrate(&UserSetting{})
	defaultUserSetting := func(k pb.UserSettingKey, v string) *UserSetting {
		return &UserSetting{
			UserId: "",
			Key:    k,
			Value:  v,
		}
	}
	insertIfNotExist := func(setting *UserSetting) {
		if xorm.RecordNotFound(tx.Where("userId = ? and `key` = ?", setting.UserId, setting.Key).First(&UserSetting{}).Error) {
			tx.Create(setting)
		}
	}
	insertIfNotExist(defaultUserSetting(pb.UserSettingKey_HowToAddFriend, "need_confirm"))
	insertIfNotExist(defaultUserSetting(pb.UserSettingKey_HowToAddFriend_NeedAnswerQuestionCorrectly_Question, "你的名字是？"))
	insertIfNotExist(defaultUserSetting(pb.UserSettingKey_HowToAddFriend_NeedAnswerQuestionCorrectly_Answer, "xxim"))
}
