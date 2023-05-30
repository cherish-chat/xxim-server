package store

import (
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

type xDatabase struct {
	sqlite *gorm.DB
	Config xConfig
}

var Database *xDatabase

func init() {
	Database = &xDatabase{}
	db, err := gorm.Open(sqlite.Open("./imcloudx.db"), &gorm.Config{})
	if err != nil {
		logx.Errorf("init sqlite error: %v", err)
		os.Exit(1)
	}
	Database.sqlite = db
	db.AutoMigrate(
		&ConfigModel{},
	)
}
