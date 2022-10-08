package svc

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/xx/internal/config"
	"github.com/cherish-chat/xxim-server/common/dbmodel"
	"github.com/cherish-chat/xxim-server/common/xredis"
	"github.com/go-redis/redis/v8"
	"github.com/qiniu/qmgo"
	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	Config         config.Config
	mongoClient    *qmgo.Client
	mongoDatabase  *qmgo.Database
	msgCollection  *qmgo.Collection
	userCollection *qmgo.Collection
	redis          redis.UniversalClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}

func (c *ServiceContext) MongoClient() *qmgo.Client {
	if c.mongoClient == nil {
		var err error
		c.mongoClient, err = qmgo.NewClient(context.Background(), &qmgo.Config{Uri: c.Config.Mongo.Uri})
		if err != nil {
			logx.Errorf("connect to mongo failed, err: %v", err)
			panic(err)
		}
	}
	return c.mongoClient
}

func (c *ServiceContext) MongoDatabase() *qmgo.Database {
	if c.mongoDatabase == nil {
		c.mongoDatabase = c.MongoClient().Database(c.Config.Mongo.Database)
	}
	return c.mongoDatabase
}

func (c *ServiceContext) MsgCollection() *qmgo.Collection {
	if c.msgCollection == nil {
		c.msgCollection = c.MongoDatabase().Collection(c.Config.Mongo.Collections.Msg)
		dbmodel.InitMsg(c.msgCollection)
	}
	return c.msgCollection
}

func (c *ServiceContext) UserCollection() *qmgo.Collection {
	if c.userCollection == nil {
		c.userCollection = c.MongoDatabase().Collection(c.Config.Mongo.Collections.User)
		dbmodel.InitUser(c.userCollection)
	}
	return c.userCollection
}

func (c *ServiceContext) Redis() redis.UniversalClient {
	if c.redis == nil {
		c.redis = xredis.GetClient(c.Config.Redis)
	}
	return c.redis
}
