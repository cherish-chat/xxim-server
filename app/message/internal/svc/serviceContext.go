package svc

import (
	"github.com/cherish-chat/xxim-server/app/message/internal/config"
	"github.com/cherish-chat/xxim-server/app/message/messagemodel"
	"github.com/cherish-chat/xxim-server/common/xcache"
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/qiniu/qmgo"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config            config.Config
	Redis             *redis.Redis
	Mysql             *gorm.DB
	MessageCollection *qmgo.QmgoClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	s := &ServiceContext{
		Config:            c,
		Redis:             xcache.MustNewRedis(c.RedisConf),
		Mysql:             xorm.MustNewMysql(c.MysqlConf),
		MessageCollection: xmgo.MustNewMongoCollection(c.Message.MongoCollection, &messagemodel.Message{}),
	}
	messagemodel.InitMessageModel(s.MessageCollection, s.Redis)
	return s
}
