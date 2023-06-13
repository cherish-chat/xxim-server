package svc

import (
	"github.com/cherish-chat/xxim-server/app/notice/internal/config"
	"github.com/cherish-chat/xxim-server/app/notice/noticemodel"
	"github.com/cherish-chat/xxim-server/common/xcache"
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/qiniu/qmgo"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config           config.Config
	Redis            *redis.Redis
	Mysql            *gorm.DB
	NoticeCollection *qmgo.QmgoClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	s := &ServiceContext{
		Config:           c,
		Redis:            xcache.MustNewRedis(c.RedisConf),
		Mysql:            xorm.MustNewMysql(c.MysqlConf),
		NoticeCollection: xmgo.MustNewMongoCollection(c.Notice.MongoCollection, &noticemodel.Notice{}),
	}
	noticemodel.InitNoticeModel(s.NoticeCollection, s.Redis)
	return s
}
