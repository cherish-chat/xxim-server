package appmgmtmodel

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"gorm.io/gorm"
)

// Link // 外部链接
type Link struct {
	Id         string `json:"id" gorm:"column:id;type:varchar(32);primary_key;not null"`
	Sort       int32  `json:"sort" gorm:"column:sort;type:int(11);not null;default:0"`
	Name       string `json:"name" gorm:"column:name;type:varchar(255);not null"`
	Url        string `json:"url" gorm:"column:url;type:varchar(255);not null"`
	Icon       string `json:"icon" gorm:"column:icon;type:varchar(255);not null"`
	IsEnable   bool   `json:"isEnable" gorm:"column:isEnable;type:tinyint(1);not null;default:1"`
	CreateTime int64  `json:"createTime" gorm:"column:createTime;type:bigint(20);not null;default:0"`
}

func (m *Link) TableName() string {
	return APPMGR_TABLE_PREFIX + "link"
}

func (m *Link) Insert(tx *gorm.DB) error {
	return tx.Create(m).Error
}

func (m *Link) ToPB() *pb.AppMgmtLink {
	return &pb.AppMgmtLink{
		Id:           m.Id,
		Sort:         m.Sort,
		Name:         m.Name,
		Url:          m.Url,
		Icon:         m.Icon,
		IsEnable:     m.IsEnable,
		CreatedAt:    m.CreateTime,
		CreatedAtStr: utils.TimeFormat(m.CreateTime),
	}
}
