package store

type ConfigModel struct {
	K string `gorm:"column:k;type:varchar(255);primary_key;not null"`
	V string `gorm:"column:v;type:varchar(255);not null"`
}

// TableName 表名
func (m *ConfigModel) TableName() string {
	return "config"
}

type xConfig struct {
}

func (x xConfig) FindByK(k string) string {
	var model ConfigModel
	Database.sqlite.Table("config").Where("k = ?", k).First(&model)
	return model.V
}

func (x xConfig) Save(k string, v string) {
	var model ConfigModel
	model.K = k
	model.V = v
	Database.sqlite.Table("config").Save(&model)
}
