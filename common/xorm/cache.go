package xorm

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type HashKv struct {
	Key string `json:"key" gorm:"column:key;type:varchar(255);not null;index:idx_key_hk,unique;"`
	HK  string `json:"hk" gorm:"column:hk;type:varchar(255);not null;index:idx_key_hk,unique;"`
	V   string `json:"v" gorm:"column:v;type:varchar(255);not null;default:'';"`
}

func (m *HashKv) TableName() string {
	return "cache_hash"
}

func MHSet(tx *gorm.DB, kvs ...HashKv) error {
	// 批量 upsert
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}, {Name: "hk"}},
		DoUpdates: clause.AssignmentColumns([]string{"v"}),
	}).Create(kvs).Error
}

func HMGet(tx *gorm.DB, key string, hks ...string) ([]*HashKv, error) {
	whereBuilder := tx.Model(&HashKv{})
	whereBuilder = whereBuilder.Where("`key` = ?", key)
	if len(hks) > 0 {
		whereBuilder = whereBuilder.Where("hk in (?)", hks)
	}
	var kvs []*HashKv
	err := whereBuilder.Find(&kvs).Error
	return kvs, err
}

func MHGet(tx *gorm.DB, kvs ...HashKv) ([]*HashKv, error) {
	whereBuilder := tx.Model(&HashKv{})
	for _, kv := range kvs {
		whereBuilder = whereBuilder.Or("`key` = ? and hk = ?", kv.Key, kv.HK)
	}
	var kvs2 []*HashKv
	err := whereBuilder.Find(&kvs2).Error
	return kvs2, err
}
