package usermodel

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"gorm.io/gorm"
)

type InvitationCode struct {
	Code        string `gorm:"column:code;type:varchar(32);not null;primary_key;" json:"code"`
	Remark      string `gorm:"column:remark;type:varchar(255);not null" json:"remark"`
	Creator     string `gorm:"column:creator;type:varchar(32);not null" json:"creator"`
	CreatorType int32  `gorm:"column:creatorType;type:tinyint(1);not null" json:"creatorType"` // 0:管理员 1:用户
	IsEnable    bool   `gorm:"column:isEnable;type:tinyint(1);not null" json:"isEnable"`
	// 默认会话模式 0:添加所有预设会话 1:只添加一个会话(轮询) 2:只添加一个会话(随机) 3:不添加会话
	DefaultConvMode int32 `gorm:"column:defaultConvMode;type:tinyint(1);not null" json:"defaultConvMode"`
	CreateTime      int64 `gorm:"column:createTime;type:bigint(20);not null" json:"createTime"`
}

func (m *InvitationCode) TableName() string {
	return TABLE_PREFIX + "invitation_code"
}

func (m *InvitationCode) Insert(tx *gorm.DB) error {
	return tx.Create(m).Error
}

func (m *InvitationCode) ToPB() *pb.UserInvitationCode {
	return &pb.UserInvitationCode{
		Code:            m.Code,
		Remark:          m.Remark,
		Creator:         m.Creator,
		CreatorType:     m.CreatorType,
		IsEnable:        m.IsEnable,
		DefaultConvMode: m.DefaultConvMode,
		CreateTime:      m.CreateTime,
		CreatedAt:       m.CreateTime,
		CreatedAtStr:    utils.TimeFormat(m.CreateTime),
	}
}
