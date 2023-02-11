package appmgmtmodel

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"gorm.io/gorm"
)

type Vpn struct {
	Id        string `gorm:"column:id;type:varchar(32);primary_key;not null" json:"id"`
	Name      string `gorm:"column:name;type:varchar(255);not null;default:'';" json:"name"`
	Platform  string `gorm:"column:platform;type:varchar(16);not null;default:'';index;" json:"platform"`
	Type      string `gorm:"column:type;type:varchar(16);not null;default:'';" json:"type"`
	Ip        string `gorm:"column:ip;type:varchar(32);not null;default:'';" json:"ip"`
	Port      int    `gorm:"column:port;type:int(11);not null;default:0;" json:"port"`
	Username  string `gorm:"column:username;type:varchar(255);not null;default:'';" json:"username"`
	Password  string `gorm:"column:password;type:varchar(255);not null;default:'';" json:"password"`
	SecretKey string `gorm:"column:secretKey;type:varchar(255);not null;default:'';" json:"secretKey"`

	CreateTime int64 `gorm:"column:createTime;type:bigint(20);not null;default:0;" json:"createTime"`
}

func (m *Vpn) TableName() string {
	return APPMGR_TABLE_PREFIX + "vpn"
}

func (m *Vpn) Insert(mysql *gorm.DB) error {
	return mysql.Create(m).Error
}

func (m *Vpn) ToPB() *pb.AppMgmtVpn {
	return &pb.AppMgmtVpn{
		Id:           m.Id,
		Name:         m.Name,
		Platform:     m.Platform,
		Type:         m.Type,
		Ip:           m.Ip,
		Port:         int32(m.Port),
		Username:     m.Username,
		Password:     m.Password,
		SecretKey:    m.SecretKey,
		CreatedAt:    m.CreateTime,
		CreatedAtStr: utils.TimeFormat(m.CreateTime),
	}
}
