package xmgo

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/utils"
	opts "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

type TestModel struct {
	Id string `bson:"_id"`
}

func TestErrIsDup(t *testing.T) {
	/*
	   Uri: "mongodb://172.88.10.41:27017/?replicaSet=rs0"
	   Database: "xxim_local"
	   Collection: "conversation_member"
	   MaxPoolSize: 100
	   MinPoolSize: 1
	   SocketTimeoutMS: 10000
	   ReadPreference: 2
	*/
	collection := MustNewMongoCollection(MongoCollectionConf{
		Uri:             "mongodb://172.88.10.41:27017/?replicaSet=rs0",
		Database:        "xxim_test",
		Collection:      "TestErrIsDup",
		MaxPoolSize:     0,
		MinPoolSize:     0,
		SocketTimeoutMS: 0,
		ReadPreference:  0,
	}, &TestModel{})
	{
		_, err := collection.InsertOne(context.Background(), &TestModel{Id: "1"})
		t.Logf("insert one error: %v", err)
	}
	{
		_, err := collection.InsertOne(context.Background(), &TestModel{Id: "1"})
		t.Logf("insert one error: %v", err)
	}
	{
		_, err := collection.InsertMany(context.Background(), []interface{}{&TestModel{Id: "1"}, &TestModel{Id: "2"}}, opts.InsertManyOptions{
			InsertManyOptions: options.InsertMany().SetOrdered(false),
		})
		t.Logf("insert many error: %v", err)
	}
}

type TestModel2 struct {
	K1  string
	K2  string
	K3  string
	K4  string
	K5  string
	K6  string
	K7  string
	K8  string
	K9  string
	K10 string
}

func TestBatchInsertMany(t *testing.T) {
	collection := MustNewMongoCollection(MongoCollectionConf{
		Uri:             "mongodb://172.88.10.41:27017/?replicaSet=rs0",
		Database:        "xxim_test",
		Collection:      "TestBatchInsertMany",
		MaxPoolSize:     0,
		MinPoolSize:     0,
		SocketTimeoutMS: 0,
		ReadPreference:  0,
	}, &TestModel2{})
	collection.DropCollection(context.Background())
	// 造数据
	var data []*TestModel2
	for i := 0; i < 1000; i++ {
		data = append(data, &TestModel2{
			K1:  utils.Snowflake.String(),
			K2:  utils.Snowflake.String(),
			K3:  utils.Snowflake.String(),
			K4:  utils.Snowflake.String(),
			K5:  utils.Snowflake.String(),
			K6:  utils.Snowflake.String(),
			K7:  utils.Snowflake.String(),
			K8:  utils.Snowflake.String(),
			K9:  utils.Snowflake.String(),
			K10: utils.Snowflake.String(),
		})
	}
	// 插入
	startTime := time.Now()
	err := BatchInsertMany(collection, context.Background(), data, 1000)
	endTime := time.Now()
	if err != nil {
		t.Fatalf("batch insert many error: %v", err)
	}
	//打印耗时
	// batchSize=500: 100000条数据，耗时：5.298955584s
	// batchSize=1000: 100000条数据，耗时：4.886134042s
	// batchSize=1500: 100000条数据，耗时：4.639113292s
	// batchSize=1500: 100000条数据，耗时：4.857033375s
	// batchSize=3000: 100000条数据，耗时：5.565197667s

	// batchSize=1000: 1000条数据，耗时：648.023542ms
	// batchSize=1000: 20000条数据，耗时：1.262419375s
	// batchSize=1500: 20000条数据，耗时：1.46495275s
	t.Logf("batch insert many cost: %v", endTime.Sub(startTime))
}
