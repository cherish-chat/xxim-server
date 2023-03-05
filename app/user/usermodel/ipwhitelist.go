package usermodel

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"gorm.io/gorm"
)

type IpList struct {
	Id         string `gorm:"column:id;type:varchar(32);primary_key;not null" json:"id"`
	Platform   string `gorm:"column:platform;type:varchar(32);not null;default:'';" json:"platform"` // 针对平台
	StartIp    string `gorm:"column:startIp;not null;default:'';comment:'起始IP';index;"`
	EndIp      string `gorm:"column:endIp;not null;default:'';comment:'结束IP';index;"`
	Remark     string `gorm:"column:remark;not null;default:'';comment:'备注';"`
	UserId     string `gorm:"column:userId;not null;default:'';comment:'用户ID';"` // 如果为空，则表示针对所有用户
	IsEnable   bool   `gorm:"column:isEnable;not null;default:0;comment:'是否启用';index;"`
	CreateTime int64  `gorm:"column:createTime;not null;default:0;comment:'创建时间';index;"`
}

type IpWhiteList struct {
	IpList
}

func (m *IpWhiteList) TableName() string {
	return TABLE_PREFIX + "ip_whitelist"
}

type IpBlackList struct {
	IpList
}

func (m *IpBlackList) TableName() string {
	return TABLE_PREFIX + "ip_blacklist"
}

func (m *IpBlackList) Insert(tx *gorm.DB) error {
	return tx.Create(m).Error
}

func (m *IpList) ToPB() *pb.UserIpList {
	return &pb.UserIpList{
		Id:           m.Id,
		Platform:     m.Platform,
		StartIp:      m.StartIp,
		EndIp:        m.EndIp,
		Remark:       m.Remark,
		UserId:       m.UserId,
		IsEnable:     m.IsEnable,
		CreateTime:   m.CreateTime,
		CreatedAt:    m.CreateTime,
		CreatedAtStr: utils.TimeFormat(m.CreateTime),
	}
}

func (m *IpWhiteList) Insert(tx *gorm.DB) error {
	return tx.Create(m).Error
}
