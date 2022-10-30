package xmgo

import (
	"context"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	opts "go.mongodb.org/mongo-driver/mongo/options"
)

type MHSetKv struct {
	Key string
	HK  string
	V   any
}

func (m *MHSetKv) CollectionName() string {
	return "cache_hash"
}

func (m *MHSetKv) Indexes(c *qmgo.Collection) error {
	_ = c.CreateIndexes(context.Background(), []options.IndexModel{{
		Key:          []string{"key", "hk"},
		IndexOptions: opts.Index().SetUnique(true),
	}, {
		Key: []string{"key"},
	}})
	return nil
}

func MHSet(c *qmgo.Collection, ctx context.Context, kvs ...MHSetKv) error {
	_, err := c.InsertMany(ctx, kvs, options.InsertManyOptions{InsertManyOptions: opts.InsertMany().SetOrdered(false)})
	return err
}

func MHGet(c *qmgo.Collection, ctx context.Context, kvs ...MHSetKv) ([]*MHSetKv, error) {
	var err error
	var ors []bson.M
	for _, kv := range kvs {
		ors = append(ors, bson.M{"key": kv.Key, "hk": kv.HK})
	}
	var results []*MHSetKv
	err = c.Find(ctx, bson.M{"$or": ors}).All(&results)
	return results, err
}

func HMGet(c *qmgo.Collection, ctx context.Context, key string, hks ...string) ([]*MHSetKv, error) {
	var err error
	var results []*MHSetKv
	err = c.Find(ctx, bson.M{
		"key": key,
		"hk":  bson.M{"$in": hks},
	}).All(&results)
	return results, err
}

func HGetAll(c *qmgo.Collection, ctx context.Context, key string) ([]*MHSetKv, error) {
	var err error
	var results []*MHSetKv
	err = c.Find(ctx, bson.M{
		"key": key,
	}).All(&results)
	return results, err
}
