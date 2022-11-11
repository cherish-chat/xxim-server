package groupmodel

import (
	"context"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	opts "go.mongodb.org/mongo-driver/mongo/options"
)

type (
	GroupMember struct {
		// 群id
		GroupId string `bson:"groupId" json:"groupId"`
		// 用户id
		UserId string `bson:"userId" json:"userId"`
		// 加入时间
		CreateTime int64 `bson:"createTime" json:"createTime"`
	}
)

func (m *GroupMember) CollectionName() string {
	return "group_member"
}

func (m *GroupMember) Indexes(c *qmgo.Collection) error {
	_ = c.CreateIndexes(context.Background(), []options.IndexModel{{
		Key:          []string{"groupId", "userId"},
		IndexOptions: opts.Index().SetUnique(true),
	}, {
		Key: []string{"groupId"},
	}, {
		Key: []string{"userId"},
	}, {
		Key: []string{"createTime"},
	}})
	return nil
}
