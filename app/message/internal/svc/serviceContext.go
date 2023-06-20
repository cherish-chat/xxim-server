package svc

import (
	"github.com/cherish-chat/xxim-server/app/message/internal/config"
	"github.com/cherish-chat/xxim-server/app/message/messagemodel"
	"github.com/cherish-chat/xxim-server/app/message/noticemodel"
	"github.com/cherish-chat/xxim-server/common/xcache"
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"github.com/cherish-chat/xxim-server/common/xmq"
	"github.com/qiniu/qmgo"
	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type ServiceContext struct {
	Config config.Config
	Redis  *redis.Redis

	MQ xmq.MQ

	BroadcastNoticeCollection           *qmgo.QmgoClient
	SubscriptionNoticeCollection        *qmgo.QmgoClient
	SubscriptionNoticeContentCollection *qmgo.QmgoClient
	MessageCollection                   *qmgo.QmgoClient

	SendMsgTokenLimiter *limit.TokenLimiter
}

func NewServiceContext(c config.Config) *ServiceContext {
	s := &ServiceContext{
		Config:                              c,
		Redis:                               xcache.MustNewRedis(c.RedisConf),
		BroadcastNoticeCollection:           xmgo.MustNewMongoCollection(c.MongoCollection.BroadcastNotice, &noticemodel.BroadcastNotice{}),
		SubscriptionNoticeCollection:        xmgo.MustNewMongoCollection(c.MongoCollection.SubscriptionNotice, &noticemodel.SubscriptionNotice{}),
		SubscriptionNoticeContentCollection: xmgo.MustNewMongoCollection(c.MongoCollection.SubscriptionNoticeContent, &noticemodel.SubscriptionNoticeContent{}),
		MessageCollection:                   xmgo.MustNewMongoCollection(c.MongoCollection.Message, &messagemodel.Message{}),
	}

	s.SendMsgTokenLimiter = limit.NewTokenLimiter(
		c.SendMsgLimiter.Rate, c.SendMsgLimiter.Burst, s.Redis, c.SendMsgLimiter.Key,
	)

	messagemodel.InitRedisSeq(s.Redis)

	s.MQ = xmq.NewAsynq(s.Config.RedisConf, 2, s.Config.Log.Level)
	return s
}
