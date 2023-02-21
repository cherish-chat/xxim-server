package usermodel

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
	"testing"
)

func TestGetUsersByIds(t *testing.T) {
	rc := getRc()
	tx := getTx()
	GetUsersByIds(context.Background(), rc, tx, []string{"showurl"})
}

func getTx() *gorm.DB {
	return xorm.NewClient(xorm.MysqlConfig{
		Addr:         "root:123456@tcp(127.0.0.1:6805)/xxim?charset=utf8mb4&parseTime=True&loc=Local&timeout=20s&readTimeout=20s&writeTimeout=20s",
		MaxIdleConns: 10,
		MaxOpenConns: 100,
		LogLevel:     "INFO",
	})
}

func getRc() *redis.Redis {
	conf := redis.RedisConf{
		Host: "127.0.0.1:6803",
		Type: "node",
		Pass: "123456",
		Tls:  false,
	}
	return conf.NewRedis()
}
