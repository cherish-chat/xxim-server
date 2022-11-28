package xorm

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// Transaction mysql 事务
func Transaction(tx *gorm.DB, fs ...func(tx *gorm.DB) error) error {
	return tx.Transaction(func(tx *gorm.DB) error {
		var err error
		for _, f := range fs {
			if f != nil {
				err = f(tx)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}

// Upsert mysql upsert
// updateKeys 更新的字段 如果为空则更新所有字段
// findKeys 查询的字段 不可以为空
func Upsert(tx *gorm.DB, model interface{}, updateKeys, findKeys []string) error {
	if len(findKeys) == 0 {
		panic("findKeys is empty")
	}
	var columns []clause.Column
	for _, key := range findKeys {
		columns = append(columns, clause.Column{Name: key})
	}
	var err error
	if len(updateKeys) > 0 {
		err = tx.Model(model).Table(model.(schema.Tabler).TableName()).Clauses(clause.OnConflict{
			Columns:   columns,
			DoUpdates: clause.AssignmentColumns(updateKeys),
		}).Create(model).Error
	} else {
		err = tx.Model(model).Table(model.(schema.Tabler).TableName()).Clauses(clause.OnConflict{
			Columns:   columns,
			UpdateAll: true,
		}).Create(model).Error
	}
	return err
}

// DetailByWhere 查询单条记录
func DetailByWhere(tx *gorm.DB, model interface{}, wheres ...GormWhere) error {
	tableName := model.(schema.Tabler).TableName()
	tx = tx.Table(tableName)
	for _, where := range wheres {
		tx = tx.Where(where.Where, where.args...)
	}
	err := tx.First(model).Error
	if err != nil {
		// 表不存在
		if TableNotFound(err) {
			// 创建表
			_ = tx.Table(tableName).AutoMigrate(model)
		}
		return err
	}
	return nil
}

// Count 获取数量
func Count(tx *gorm.DB, model interface{}, where string, args ...interface{}) (int64, error) {
	var total int64
	err := tx.Model(model).Where(where, args...).Count(&total).Error
	if TableNotFound(err) {
		_ = tx.AutoMigrate(model)
		err = tx.Model(model).Where(where, args...).Count(&total).Error
	}
	return total, err
}

// InsertOne 插入一条记录
func InsertOne(tx *gorm.DB, model interface{}) error {
	tableName := model.(schema.Tabler).TableName()
	err := tx.Table(tableName).Create(model).Error
	if err != nil {
		// 表不存在
		if TableNotFound(err) {
			// 创建表
			err = tx.Table(tableName).AutoMigrate(model)
			if err != nil {
				return err
			} else {
				// 创建记录
				return tx.Table(tableName).Create(model).Error
			}
		} else {
			return err
		}
	}
	return nil
}

// MustUpdate 必须更新 如果没有记录则报错
func MustUpdate(
	tx *gorm.DB, model interface{},
	updateMap map[string]interface{}, wheres ...GormWhere,
) error {
	tableName := model.(schema.Tabler).TableName()
	db := tx.Model(model).Table(tableName)
	for _, where := range wheres {
		db = db.Where(where.Where, where.args...)
	}
	updates := db.Updates(updateMap)
	err := updates.Error
	if err != nil {
		return err
	}
	if updates.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Update 更新
func Update(
	tx *gorm.DB, model interface{},
	updateMap map[string]interface{}, wheres ...GormWhere,
) error {
	err := MustUpdate(tx, model, updateMap, wheres...)
	if err != nil {
		if RecordNotFound(err) {
			return nil
		}
		return err
	}
	return nil
}

// InsertMany 插入多条记录
func InsertMany(tx *gorm.DB, model interface{}, models interface{}) error {
	tableName := model.(schema.Tabler).TableName()
	err := tx.Table(tableName).Create(models).Error
	if err != nil {
		// 表不存在
		if TableNotFound(err) {
			// 创建表
			err = tx.Table(tableName).AutoMigrate(model)
			if err != nil {
				return err
			} else {
				// 创建记录
				return tx.Table(tableName).Create(models).Error
			}
		} else {
			return err
		}
	}
	return nil
}

// ListWithPaging 分页查询
func ListWithPaging(
	tx *gorm.DB,
	models interface{},
	model interface{},
	no int, size int,
	where string, args ...interface{}) (int64, error) {
	tableName := model.(schema.Tabler).TableName()
	var count int64
	db := tx.Table(tableName).Where(where, args...)
	db.Count(&count)
	return count, Paging(db, no, size).Find(models).Error
}

// ListWithPagingOrder 分页查询
func ListWithPagingOrder(
	tx *gorm.DB,
	models interface{},
	model interface{},
	no int, size int,
	order string,
	where string, args ...interface{}) (int64, error) {
	tableName := model.(schema.Tabler).TableName()
	var count int64
	db := tx.Table(tableName).Where(where, args...)
	db.Count(&count)
	return count, Paging(db.Order(order), no, size).Find(models).Error
}

// Paging 分页
func Paging(tx *gorm.DB, no int, size int) *gorm.DB {
	if size != 0 {
		return tx.Offset((no - 1) * size).Limit(size)
	}
	return tx
}

// Pluck 获取某个字段的值
func Pluck(tx *gorm.DB, column string, values interface{}, model interface{}, limit int, where string, args ...interface{}) error {
	return tx.Model(model).Where(where, args...).Limit(limit).Pluck(column, values).Error
}
