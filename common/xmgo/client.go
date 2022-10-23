package xmgo

import (
	"context"
	"github.com/qiniu/qmgo"
	"github.com/zeromicro/go-zero/core/logx"
	"sync"
)

type Config struct {
	Uri      string
	Database string
}

type Client struct {
	*qmgo.Client
	config        Config
	database      *qmgo.Database
	CollectionMap sync.Map
}

func NewClient(
	cfg Config,
) *Client {
	mongoClient, err := qmgo.NewClient(context.Background(), &qmgo.Config{
		Uri: cfg.Uri,
	})
	if err != nil {
		panic(err)
	}
	err = mongoClient.Ping(1)
	if err != nil {
		logx.Errorf("ping mongo err: %v", err)
		panic(err)
	}
	return &Client{
		Client:        mongoClient,
		config:        cfg,
		database:      mongoClient.Database(cfg.Database),
		CollectionMap: sync.Map{},
	}
}

func (c *Client) Collection(i ICollection) *qmgo.Collection {
	if v, ok := c.CollectionMap.Load(i.CollectionName()); ok {
		return v.(*qmgo.Collection)
	} else {
		collection := c.database.Collection(i.CollectionName())
		err := i.Indexes(collection)
		if err != nil {
			logx.Errorf("create collection index err: %v", err)
			panic(err)
		}
		c.CollectionMap.Store(i.CollectionName(), collection)
		return collection
	}
}
