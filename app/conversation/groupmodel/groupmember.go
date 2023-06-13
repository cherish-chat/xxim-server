package groupmodel

import (
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
	"os"
)

type GroupMember struct {
	// GroupId 群组ID
	GroupId int `gorm:"column:groupId;type:bigint(20);not null;index:idx_groupId;primary_key;"`
	// MemberUserId 群成员ID
	MemberUserId int `gorm:"column:memberUserId;type:bigint(20);not null;index:idx_memberUserId;primary_key;"`
	// JoinTime 加入时间
	JoinTime int `gorm:"column:joinTime;type:int(11);not null;"`
	// JoinSource 加入来源
	JoinSource string `gorm:"column:joinSource;type:varchar(32);not null;default:'';"`
}

func (m *GroupMember) TableName() string {
	return "group_member"
}

type xGroupMemberModel struct {
	db *gorm.DB
	rc *redis.Redis
}

var GroupMemberModel *xGroupMemberModel

func InitGroupMemberModel(db *gorm.DB, rc *redis.Redis) {
	GroupMemberModel = &xGroupMemberModel{
		db: db,
		rc: rc,
	}
	err := db.AutoMigrate(&GroupMember{})
	if err != nil {
		logx.Errorf("db.AutoMigrate(&GroupMember{}) error: %v", err)
		os.Exit(1)
		return
	}
}
