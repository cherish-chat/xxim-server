package appmgmtmodel

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"gorm.io/gorm"
)

type Version struct {
	Id string `gorm:"column:id;type:varchar(255);primary_key;not null" json:"id"`
	// 版本号
	Version string `gorm:"column:version;type:varchar(255);not null;default:'';index;" json:"version"`
	// 更新平台
	Platform string `gorm:"column:platform;type:varchar(255);not null;default:'';index;" json:"platform"`
	// 更新类型
	Type int8 `gorm:"column:type;type:tinyint(1);default:0;" json:"type"` // 0: 不提示 1: 提示 2: 强制
	// 更新内容
	Content string `gorm:"column:content;type:text;" json:"content"`
	// 下载地址
	DownloadUrl string `gorm:"column:downloadUrl;type:varchar(255);not null;default:'';" json:"downloadUrl"`

	CreateTime int64 `gorm:"column:createTime;type:bigint(20);not null;default:0;" json:"createTime"`
}

func (m *Version) TableName() string {
	return APPMGR_TABLE_PREFIX + "version"
}

func (m *Version) Insert(mysql *gorm.DB) error {
	return mysql.Create(m).Error
}

func (m *Version) ToPB() *pb.AppMgmtVersion {
	return &pb.AppMgmtVersion{
		Id:           m.Id,
		Version:      m.Version,
		Platform:     m.Platform,
		Type:         int32(m.Type),
		Content:      m.Content,
		DownloadUrl:  m.DownloadUrl,
		CreatedAt:    m.CreateTime,
		CreatedAtStr: utils.TimeFormat(m.CreateTime),
	}
}
