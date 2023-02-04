package appmgmtmodel

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"strings"
)

type Config struct {
	Group string `gorm:"column:group;type:varchar(255);"`
	K     string `gorm:"column:k;type:varchar(255);primary_key;not null" json:"k"`
	V     string `gorm:"column:v;type:varchar(255);"`
	Type  string `gorm:"column:type;type:varchar(255);"` // 1: string, 2: number, 3: bool, 4: json, 5: stringArray
	Name  string `gorm:"column:name;type:varchar(255);"`
	// 作用平台
	ScopePlatforms string `gorm:"column:scopePlatforms;type:varchar(255);"` // 空: 所有平台, ios android windows mac linux web system
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
	}
}
