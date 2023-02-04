package mgmtmodel

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"gorm.io/gorm"
	"time"
)

type MSIPWhitelist struct {
	Id         string `gorm:"column:id;primarykey;comment:'主键'"`
	StartIp    string `gorm:"column:startIp;not null;default:'';comment:'起始IP';"`
	EndIp      string `gorm:"column:endIp;not null;default:'';comment:'结束IP';"`
	Remark     string `gorm:"column:remark;not null;default:'';comment:'备注';"`
	UserId     string `gorm:"column:userId;not null;default:'';comment:'用户ID';"` // 如果为空，则表示针对所有用户
	IsEnable   bool   `gorm:"column:isEnable;not null;default:0;comment:'是否启用';index;"`
	CreateTime int64  `gorm:"column:createTime;not null;default:0;comment:'创建时间';index;"`
}

func (m *MSIPWhitelist) TableName() string {
	return MGMT_TABLE_PREFIX + "ms_ip_whitelist"
}

func (m *MSIPWhitelist) ToPB() *pb.MSIpWhiteList {
	return &pb.MSIpWhiteList{
		Id:           m.Id,
		StartIp:      m.StartIp,
		EndIp:        m.EndIp,
		Remark:       m.Remark,
		UserId:       m.UserId,
		IsEnable:     m.IsEnable,
		CreatedAt:    m.CreateTime,
		CreatedAtStr: utils.TimeFormat(m.CreateTime),
	}
}

func initMSIPWhitelist(tx *gorm.DB) {
	tx.Create(&MSIPWhitelist{
		Id:         "1",
		StartIp:    "0.0.0.0",
		EndIp:      "255.255.255.0",
		Remark:     "默认全部IP",
		UserId:     "",
		IsEnable:   true,
		CreateTime: time.Now().UnixMilli(),
	})
}
