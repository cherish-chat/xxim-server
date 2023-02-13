package usermodel

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"gorm.io/gorm"
)

type DefaultConv struct {
	Id string `gorm:"column:id;type:varchar(32);primary_key;not null" json:"id"`
	// 会话类型
	ConvType int32 `gorm:"column:convType;type:tinyint(1);not null" json:"convType"` // 0:好友 1:群聊
	// 类型
	FilterType     int32  `gorm:"column:filterType;type:tinyint(1);not null" json:"filterType"`                     // 0:所有注册用户 1:使用邀请码注册用户
	InvitationCode string `gorm:"column:invitationCode;type:varchar(32);not null;default:0;" json:"invitationCode"` // 空表示所有邀请码用户

	ConvId string `gorm:"column:convId;type:varchar(32);not null" json:"convId"` // 会话id

	CreateTime int64  `gorm:"column:createTime;type:bigint(20);not null" json:"createTime"`
	TextMsg    string `gorm:"column:textMsg;type:varchar(255);not null;default:''" json:"textMsg"` // 文本消息
}

func (m *DefaultConv) TableName() string {
	return TABLE_PREFIX + "default_conv"
}

func (m *DefaultConv) Insert(tx *gorm.DB) error {
	return tx.Create(m).Error
}

func (m *DefaultConv) ToPB() *pb.UserDefaultConv {
	return &pb.UserDefaultConv{
		Id:             m.Id,
		ConvType:       m.ConvType,
		FilterType:     m.FilterType,
		InvitationCode: m.InvitationCode,
		ConvId:         m.ConvId,
		TextMsg:        m.TextMsg,
		CreatedAt:      m.CreateTime,
		CreatedAtStr:   utils.TimeFormat(m.CreateTime),
	}
}
