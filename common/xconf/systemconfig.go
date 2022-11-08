package xconf

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	opts "go.mongodb.org/mongo-driver/mongo/options"
)

type SystemConfig struct {
	Namespace string `bson:"namespace" json:"namespace"`
	// 服务名称
	ServiceName string `bson:"serviceName" json:"serviceName"`
	// 键
	Key string `bson:"key" json:"key"`
	// 值
	Value string `bson:"value" json:"value"`
	// 输入选项
	InputOptions bson.M `bson:"inputOptions,omitempty" json:"inputOptions,omitempty"`
}

func (m *SystemConfig) CollectionName() string {
	return "system_config"
}

func (m *SystemConfig) Indexes(c *qmgo.Collection) error {
	_ = c.CreateIndexes(context.Background(), []options.IndexModel{{
		Key:          []string{"namespace", "serviceName", "key"},
		IndexOptions: opts.Index().SetUnique(true),
	}, {
		Key: []string{"namespace", "serviceName"},
	}, {
		Key: []string{"namespace"},
	}})
	return nil
}

type SystemConfigMgr struct {
	Namespace   string
	ServiceName string
	c           *qmgo.Collection
}

func NewSystemConfigMgr(namespace string, serviceName string, c *qmgo.Collection) *SystemConfigMgr {
	return &SystemConfigMgr{Namespace: namespace, ServiceName: serviceName, c: c}
}

func (s *SystemConfigMgr) Get(key string) (value string) {
	return s.GetCtx(context.Background(), key)
}

func (s *SystemConfigMgr) GetCtx(ctx context.Context, key string) (value string) {
	var config SystemConfig
	err := s.c.Find(ctx, bson.M{
		"namespace":   s.Namespace,
		"serviceName": s.ServiceName,
		"key":         key,
	}).One(&config)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// 创建
			config = SystemConfig{
				Namespace:   s.Namespace,
				ServiceName: s.ServiceName,
				Key:         key,
				Value:       "",
			}
			_, _ = s.c.Upsert(ctx, bson.M{
				"namespace":   s.Namespace,
				"serviceName": s.ServiceName,
				"key":         key,
			}, config)
		}
		return ""
	}
	return config.Value
}

func (s *SystemConfigMgr) GetOrDefault(key string, defaultValue string) (value string) {
	ctx := context.Background()
	return s.GetOrDefaultCtx(ctx, key, defaultValue)
}

func (s *SystemConfigMgr) GetOrDefaultCtx(ctx context.Context, key string, defaultValue string) (value string) {
	var config SystemConfig
	err := s.c.Find(ctx, bson.M{
		"namespace":   s.Namespace,
		"serviceName": s.ServiceName,
		"key":         key,
	}).One(&config)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// 创建
			config = SystemConfig{
				Namespace:   s.Namespace,
				ServiceName: s.ServiceName,
				Key:         key,
				Value:       "",
			}
			_, _ = s.c.Upsert(ctx, bson.M{
				"namespace":   s.Namespace,
				"serviceName": s.ServiceName,
				"key":         key,
			}, config)
		}
		return defaultValue
	}
	return config.Value
}

func (s *SystemConfigMgr) GetSlice(key string) []string {
	return s.GetSliceCtx(context.Background(), key)
}

func (s *SystemConfigMgr) GetSliceCtx(ctx context.Context, key string) []string {
	var config SystemConfig
	err := s.c.Find(ctx, bson.M{
		"namespace":   s.Namespace,
		"serviceName": s.ServiceName,
		"key":         key,
	}).One(&config)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			// 创建
			config = SystemConfig{
				Namespace:   s.Namespace,
				ServiceName: s.ServiceName,
				Key:         key,
				Value:       "[]",
			}
			_, _ = s.c.Upsert(ctx, bson.M{
				"namespace":   s.Namespace,
				"serviceName": s.ServiceName,
				"key":         key,
			}, config)
		}
		return nil
	}
	res := make([]string, 0)
	_ = json.Unmarshal([]byte(config.Value), &res)
	return res
}
