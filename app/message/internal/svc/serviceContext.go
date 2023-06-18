package svc

import (
	"github.com/cherish-chat/xxim-server/app/message/internal/config"
	"github.com/cherish-chat/xxim-server/app/message/noticemodel"
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"github.com/cherish-chat/xxim-server/common/xmq"
	"github.com/qiniu/qmgo"
)

type ServiceContext struct {
	Config                              config.Config
	MQ                                  xmq.MQ
	BroadcastNoticeCollection           *qmgo.QmgoClient
	SubscriptionNoticeCollection        *qmgo.QmgoClient
	SubscriptionNoticeContentCollection *qmgo.QmgoClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	s := &ServiceContext{
		Config:                              c,
		BroadcastNoticeCollection:           xmgo.MustNewMongoCollection(c.MongoCollection.BroadcastNotice, &noticemodel.BroadcastNotice{}),
		SubscriptionNoticeCollection:        xmgo.MustNewMongoCollection(c.MongoCollection.SubscriptionNotice, &noticemodel.SubscriptionNotice{}),
		SubscriptionNoticeContentCollection: xmgo.MustNewMongoCollection(c.MongoCollection.SubscriptionNoticeContent, &noticemodel.SubscriptionNoticeContent{}),
	}
	s.MQ = xmq.NewAsynq(s.Config.RedisConf, 2, s.Config.Log.Level)
	return s
}
