package mgmtmodel

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"gorm.io/gorm"
	"strings"
	"time"
)

type Role struct {
	Id         string `gorm:"column:id;primarykey;comment:'主键'"`
	Name       string `gorm:"column:name;not null;default:'';comment:'角色名称''"`
	Remark     string `gorm:"column:remark;not null;default:'';comment:'备注信息'"`
	IsDisable  bool   `gorm:"column:isDisable;not null;default:0;comment:'是否禁用: 0=否, 1=是'"`
	Sort       int32  `gorm:"column:sort;not null;default:0;comment:'角色排序'"`
	MenuIds    string `gorm:"column:menuIds;not null;default:'';comment:'菜单ID集合逗号分隔'"`
	ApiPathIds string `gorm:"column:apiPathIds;not null;default:'';comment:'接口ID集合逗号分隔'"`
	CreateTime int64  `gorm:"column:createTime;not null;comment:'创建时间'"`
	UpdateTime int64  `gorm:"column:updateTime;not null;comment:'更新时间'"`
}

func (m *Role) TableName() string {
	return MGMT_TABLE_PREFIX + "role"
}

func (m *Role) ToPB() *pb.MSRole {
	return &pb.MSRole{
		Id:           m.Id,
		CreatedAt:    m.CreateTime,
		CreatedAtStr: utils.TimeFormat(m.CreateTime),
		CreatedBy:    "",
		UpdatedAt:    m.UpdateTime,
		UpdatedAtStr: utils.TimeFormat(m.UpdateTime),
		UpdatedBy:    "",
		Name:         m.Name,
		Remark:       m.Remark,
		IsDisable:    m.IsDisable,
		Sort:         m.Sort,
		ApiPathIds:   strings.Split(m.ApiPathIds, ","),
		ApiPaths:     nil,
		MenuIds:      strings.Split(m.MenuIds, ","),
		Menus:        nil,
	}
}

var defaultRoles = []*Role{
	{
		Id:         "1",
		Name:       "超级管理员",
		Remark:     "超级管理员",
		IsDisable:  false,
		Sort:       0,
		CreateTime: time.Now().UnixMilli(),
		UpdateTime: time.Now().UnixMilli(),
	},
}

func initRole(tx *gorm.DB) {
	for _, role := range defaultRoles {
		insertIfNotFound(tx, role.Id, role)
	}
}
