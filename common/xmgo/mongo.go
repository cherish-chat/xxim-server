package xmgo

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/qiniu/qmgo"
	"github.com/zeromicro/go-zero/core/logx"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"time"
)

type MongoCollectionConf struct {
	// Uri is the connection URI for the mongo instance
	// example: mongodb://$user:$pwd@$host:$port
	// example: mongodb://172.88.10.41:27017/?replicaSet=rs0
	Uri             string
	Database        string
	Collection      string
	MaxPoolSize     uint64 `json:",default=100"`
	MinPoolSize     uint64 `json:",optional"`
	SocketTimeoutMS int64  `json:",default=300000"` // default is 5 minutes
	//ReadPreference 读偏好模式
	//DefaultMode:0 默认不设置
	//PrimaryMode:1 从主节点读取数据
	//PrimaryPreferredMode:2 优先从主节点读取数据，如果主节点不可用，从从节点读取数据
	//SecondaryMode:3 从从节点读取数据
	//SecondaryPreferredMode:4 优先从从节点读取数据，如果从节点不可用，从主节点读取数据
	//NearestMode:5 从最近的成员读取数据
	ReadPreference int `json:",optional"` // default is primary
}

func MustNewMongoCollection(config MongoCollectionConf, model any) *qmgo.QmgoClient {
	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*5)
	defer cancelFunc()
	var readPreference *qmgo.ReadPref
	if config.ReadPreference > 0 && config.ReadPreference < 6 {
		readPreference = &qmgo.ReadPref{
			Mode: readpref.Mode(config.ReadPreference),
		}
	}
	qmgoClient, err := qmgo.Open(ctx, &qmgo.Config{
		Uri:              config.Uri,
		Database:         config.Database,
		Coll:             config.Collection,
		ConnectTimeoutMS: utils.AnyPtr(int64(5000)),
		MaxPoolSize:      utils.AnyPtr(config.MaxPoolSize),
		MinPoolSize:      utils.AnyPtr(config.MinPoolSize),
		SocketTimeoutMS:  utils.AnyPtr(config.SocketTimeoutMS),
		ReadPreference:   readPreference,
	})
	if err != nil {
		logx.Errorf("open mongo error: %v", err)
		os.Exit(1)
		return nil
	}
	if indexer, ok := model.(Indexer); ok {
		err = CreateIndex(ctx, qmgoClient, indexer)
		if err != nil {
			logx.Errorf("create index error: %v", err)
			os.Exit(1)
			return nil
		}
	}
	return qmgoClient
}

func CreateIndex(ctx context.Context, qmgoClient *qmgo.QmgoClient, indexer Indexer) error {
	indexes := indexer.GetIndexes()
	if len(indexes) == 0 {
		return nil
	}
	return qmgoClient.CreateIndexes(ctx, indexes)
}
