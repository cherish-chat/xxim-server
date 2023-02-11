package mgmtmodel

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
)

type LoginRecordInfo LoginInfo
type LoginRecord struct {
	Id     string `bson:"_id" json:"id" gorm:"column:id;type:char(32);primary_key"`
	UserId string `bson:"userId" json:"userId" gorm:"column:userId;type:char(32);index"`
	LoginRecordInfo
}

func (m *LoginRecord) TableName() string {
	return MGMT_TABLE_PREFIX + "login_record"
}

func (m *LoginRecord) ToPB() *pb.MSLoginRecord {
	return &pb.MSLoginRecord{
		Id:         m.Id,
		UserId:     m.UserId,
		Time:       m.Time,
		TimeStr:    utils.TimeFormat(m.Time),
		Ip:         m.Ip,
		IpCountry:  m.IpCountry,
		IpProvince: m.IpProvince,
		IpCity:     m.IpCity,
		IpISP:      m.IpISP,
		UserAgent:  m.UserAgent,
	}
}
