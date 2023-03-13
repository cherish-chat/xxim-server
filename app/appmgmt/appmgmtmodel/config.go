package appmgmtmodel

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"strings"
)

type Config struct {
	Group string `gorm:"column:group;type:varchar(255);"`
	K     string `gorm:"column:k;type:varchar(255);not null;index:uk,unique;" json:"k"`
	V     string `gorm:"column:v;type:text;"`
	Type  string `gorm:"column:type;type:varchar(255);"` // 1: string, 2: number, 3: bool, 4: json, 5: stringArray
	Name  string `gorm:"column:name;type:varchar(255);"`
	// 作用平台
	ScopePlatforms string `gorm:"column:scopePlatforms;type:varchar(255);"` // 空: 所有平台, ios android windows mac linux web system
	Options        string `gorm:"column:options;type:varchar(255);"`
	UserId         string `gorm:"column:userId;type:varchar(255);not null;default:'';index;index:uk,unique;"` // 用户id 如果为空则为全部用户的默认配置
	UpdateTime     int64  `gorm:"column:updateTime;type:bigint(20);not null;default:0;"`
}

func (m *Config) GetScopePlatforms() []string {
	if m.ScopePlatforms == "" {
		return []string{}
	}
	return strings.Split(m.ScopePlatforms, ",")
}

func (m *Config) TableName() string {
	return APPMGR_TABLE_PREFIX + "config"
}

func (m *Config) ToPB() *pb.AppMgmtConfig {
	return &pb.AppMgmtConfig{
		Group:          m.Group,
		K:              m.K,
		V:              m.V,
		Type:           m.Type,
		Name:           m.Name,
		ScopePlatforms: m.ScopePlatforms,
		Options:        m.Options,
	}
}
