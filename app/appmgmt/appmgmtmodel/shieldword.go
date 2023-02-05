package appmgmtmodel

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
)

type ShieldWord struct {
	Id         string `gorm:"column:id;type:varchar(32);primary_key;not null" json:"id"`
	Word       string `gorm:"column:word;type:varchar(255);not null;default:'';index;" json:"word"`
	CreateTime int64  `gorm:"column:createTime;type:bigint(20);not null;default:0;" json:"createTime"`
}

func (m *ShieldWord) TableName() string {
	return APPMGR_TABLE_PREFIX + "shieldword"
}

func (m *ShieldWord) ToPB() *pb.AppMgmtShieldWord {
	return &pb.AppMgmtShieldWord{
		Id:           m.Id,
		Word:         m.Word,
		CreatedAt:    m.CreateTime,
		CreatedAtStr: utils.TimeFormat(m.CreateTime),
	}
}
