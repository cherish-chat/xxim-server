package usermodel

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/utils/ip2region"
)

type LoginRecordInfo LoginInfo
type LoginRecord struct {
	Id     string `bson:"_id" json:"id" gorm:"column:id;type:char(32);primary_key"`
	UserId string `bson:"userId" json:"userId" gorm:"column:userId;type:char(32);index"`
	LoginRecordInfo
}

func (m *LoginRecord) TableName() string {
	return "login_record"
}

func (m *LoginRecord) ToPB() *pb.UserLoginRecord {
	return &pb.UserLoginRecord{
		Time:    m.Time,
		TimeStr: utils.TimeFormat(m.Time),
		Ip:      m.Ip,
		IpRegion: ip2region.Obj{
			Country:  m.IpCountry,
			Province: m.IpProvince,
			City:     m.IpCity,
			ISP:      m.IpISP,
		}.String(),
	}
}

func (m *LoginRecord) ToProto() *pb.LoginRecord {
	return &pb.LoginRecord{
		Id:          m.Id,
		UserId:      m.UserId,
		Time:        m.Time,
		TimeStr:     utils.TimeFormat(m.Time),
		Ip:          m.Ip,
		IpCountry:   m.IpCountry,
		IpProvince:  m.IpProvince,
		IpCity:      m.IpCity,
		IpISP:       m.IpISP,
		AppVersion:  m.AppVersion,
		UserAgent:   m.UserAgent,
		OsVersion:   m.OsVersion,
		Platform:    m.Platform,
		DeviceId:    m.DeviceId,
		DeviceModel: m.DeviceModel,
	}
}
