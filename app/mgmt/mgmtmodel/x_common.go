package mgmtmodel

import (
	"gorm.io/gorm"
	"strings"
)

const (
	MGMT_TABLE_PREFIX = "mgmt_"
)

func insertIfNotFound(
	tx *gorm.DB,
	id string,
	model interface{},
) {
	if tx.First(model, id).RowsAffected == 0 {
		tx.Create(model)
	}
}

func upsert(tx *gorm.DB, id string, model interface{}) {
	// delete
	tx.Where("id = ?", id).Delete(model)
	// insert
	tx.Create(model)
}

func InitData(tx *gorm.DB) {
	// restore
	initMenu(tx)
	restoreMenu(tx)
	insertMenu(tx)
	initRole(tx)
	initApiPath(tx)
	initMSIPWhitelist(tx)
	// 设置角色菜单关联
	initRoleMenu(tx)
	// 设置角色接口关联
	initRoleApiPath(tx)
	// 相册分类
	initAlbumCate(tx)
	// lua配置
	initLuaConfig(tx)
}

func initRoleMenu(tx *gorm.DB) {
	// 查询所有menu
	var menuIds []string
	tx.Model(&Menu{}).Pluck("id", &menuIds)
	// 只更新 id=1 的角色
	tx.Model(&Role{}).Where("id = ?", "1").Update("menuIds", strings.Join(menuIds, ","))
}

func initRoleApiPath(tx *gorm.DB) {
	// 查询所有apiPath
	var apiPathIds []string
	tx.Model(&ApiPath{}).Pluck("id", &apiPathIds)
	// 只更新 id=1 的角色
	tx.Model(&Role{}).Where("id = ?", "1").Update("apiPathIds", strings.Join(apiPathIds, ","))
}
