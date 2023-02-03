package mgmtmodel

import "gorm.io/gorm"

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

func InitData(tx *gorm.DB) {
	initMenu(tx)
	initRole(tx)
	initApiPath(tx)
}
