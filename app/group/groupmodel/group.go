package groupmodel

import (
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xcache"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
	"os"
)

// Group 群组 数据库模型
type Group struct {
	GroupId int `json:"groupId" gorm:"column:groupId;type:char(32);primary_key;not null"`
}

// TableName 表名
func (m *Group) TableName() string {
	return "group"
}

type xGroupModel struct {
	db *gorm.DB
	rc *redis.Redis
}

var GroupModel *xGroupModel

func InitGroupModel(db *gorm.DB, rc *redis.Redis, minGroupId int) {
	GroupModel = &xGroupModel{
		db: db,
		rc: rc,
	}
	err := db.AutoMigrate(&Group{})
	if err != nil {
		logx.Errorf("db.AutoMigrate(&Group{}) error: %v", err)
		os.Exit(1)
		return
	}
	// 查询最大的群组id
	var group Group
	err = db.Model(&Group{}).Order("groupId desc").First(&group).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 说明没有数据
			group.GroupId = minGroupId
		} else {
			logx.Errorf(`db.Model(&Group{}).Order("groupId desc").First(&group) error: %v`, err)
			os.Exit(1)
			return
		}
	}
	// 保存到redis
	err = rc.Set(xcache.IncrKeyGroupId, utils.AnyString(group.GroupId))
	if err != nil {
		logx.Errorf(`rc.Set(xcache.IncrKeyGroupId, utils.AnyString(group.GroupId)) error: %v`, err)
		os.Exit(1)
		return
	}
}

func (x *xGroupModel) GenerateGroupId() int {
	// 从redis中获取
	groupId, err := x.rc.Incr(xcache.IncrKeyGroupId)
	if err != nil {
		logx.Errorf("x.rc.Incr(xcache.IncrKeyGroupId) error: %v", err)
		return 0
	}
	return int(groupId)
}

func (x *xGroupModel) Insert(group *Group) error {
	if group.GroupId == 0 {
		group.GroupId = x.GenerateGroupId()
	}
	return x.db.Create(group).Error
}
