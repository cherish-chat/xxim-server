package mgmtmodel

import (
	"github.com/cherish-chat/xxim-server/common/xorm"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"log"
	"strconv"
)

type AutoIncrement struct {
	Service string `gorm:"column:service;type:varchar(255);primary_key" json:"service"`
	Counter uint64 `gorm:"column:counter;type:bigint(20)" json:"counter"`
}

func (m *AutoIncrement) TableName() string {
	return MGMT_TABLE_PREFIX + "auto_increment"
}

func GetId(tx *gorm.DB, tabler schema.Tabler, min int64) string {
	// 查询自增表(加行锁)
	var autoIncrement AutoIncrement
	err := tx.Set("gorm:query_option", "FOR UPDATE").Where("service = ?", tabler.TableName()).First(&autoIncrement).Error
	if err != nil {
		if xorm.RecordNotFound(err) {
			// 不存在则创建
			autoIncrement = AutoIncrement{
				Service: tabler.TableName(),
				Counter: uint64(min),
			}
			err = tx.Create(&autoIncrement).Error
			if err != nil {
				log.Fatalf("创建自增表失败: %v", err)
			}
		} else {
			log.Fatalf("查询自增表失败: %v", err)
		}
	}
	// 自增
	autoIncrement.Counter++
	err = tx.Save(&autoIncrement).Error
	if err != nil {
		log.Fatalf("自增失败: %v", err)
	}
	return strconv.FormatInt(int64(autoIncrement.Counter), 10)
}
