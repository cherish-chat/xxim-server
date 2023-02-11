package appmgmtmodel

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"gorm.io/gorm"
)

// Notice app飘屏公告
type Notice struct {
	Id         string `gorm:"column:id;type:varchar(32);primary_key;not null" json:"id"`
	Position   int8   `gorm:"column:position;type:tinyint(1);not null;default:0;index;" json:"position"` // 1: 开屏广告 2: 首页弹窗 3: 首页飘屏
	Platform   string `gorm:"column:platform;type:varchar(16);not null;default:'';index;" json:"platform"`
	Title      string `gorm:"column:title;type:varchar(255);not null;default:'';" json:"title"`
	Image      string `gorm:"column:image;type:varchar(255);not null;default:'';" json:"image"`
	Content    string `gorm:"column:content;type:text;" json:"content"`
	Sort       int32  `gorm:"column:sort;type:int(11);not null;default:0;" json:"sort"`
	IsEnable   bool   `gorm:"column:isEnable;type:tinyint(1);not null;default:0;" json:"isEnable"`
	StartTime  int64  `gorm:"column:startTime;type:bigint(20);not null;default:0;" json:"startTime"`
	EndTime    int64  `gorm:"column:endTime;type:bigint(20);not null;default:0;" json:"endTime"`
	CreateTime int64  `gorm:"column:createTime;type:bigint(20);not null;default:0;" json:"createTime"`
}

func (m *Notice) TableName() string {
	return APPMGR_TABLE_PREFIX + "notice"
}

func (m *Notice) Insert(tx *gorm.DB) error {
	return tx.Create(m).Error
}

func (m *Notice) ToPB() *pb.AppMgmtNotice {
	return &pb.AppMgmtNotice{
		Id:           m.Id,
		Position:     int32(m.Position),
		Platform:     m.Platform,
		Title:        m.Title,
		Image:        m.Image,
		Content:      m.Content,
		Sort:         m.Sort,
		IsEnable:     m.IsEnable,
		StartTime:    m.StartTime,
		EndTime:      m.EndTime,
		CreatedAt:    m.CreateTime,
		CreatedAtStr: utils.TimeFormat(m.CreateTime),
	}
}
