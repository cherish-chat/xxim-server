package immodel

import (
	"github.com/cherish-chat/xxim-server/common/utils/ip2region"
	"github.com/cherish-chat/xxim-server/common/xorm"
)

type UserConnectRecord struct {
	UserId         string        `json:"userId" bson:"userId" gorm:"column:userId;type:char(32);not null;index;comment:用户id;default:'';index:userId_deviceId;"`
	DeviceId       string        `json:"deviceId" bson:"deviceId" gorm:"column:deviceId;type:varchar(64);not null;index;comment:设备id;default:'';index:userId_deviceId;"`
	Platform       string        `json:"platform" bson:"platform" gorm:"column:platform;type:char(32);not null;comment:平台;default:''"`
	Ips            string        `json:"ips" bson:"ips" gorm:"column:ips;type:char(32);not null;comment:ip地址;default:'';"`
	IpRegion       ip2region.Obj `json:"ipRegion" bson:"ipRegion" gorm:"column:ipRegion;type:json;comment:ip地址;"`
	NetworkUsed    string        `json:"networkUsed" bson:"networkUsed" gorm:"column:networkUsed;type:char(32);not null;comment:网络类型;default:'';"`
	Headers        xorm.M        `json:"headers" bson:"headers" gorm:"column:headers;type:json;comment:请求头;"`
	PodIp          string        `json:"podIp" bson:"podIp" gorm:"column:podIp;type:char(32);not null;comment:podIp;default:'';"`
	ConnectTime    int64         `json:"connectTime" bson:"connectTime" gorm:"column:connectTime;type:bigint;not null;comment:连接时间;default:0;"`
	DisconnectTime int64         `json:"disconnectTime" bson:"disconnectTime" gorm:"column:disconnectTime;type:bigint;not null;comment:断开时间;default:0;"`
}

func (m *UserConnectRecord) TableName() string {
	return "user_connect_record"
}
