package relationmodel

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
)

type RequestAddFriend struct {
	Id string `json:"id" bson:"_id"`
	// 发起人
	FromUserId string `json:"fromUserId" bson:"fromUserId"`
	// 接收人
	ToUserId string `json:"toUserId" bson:"toUserId"`
	// 申请状态
	Status pb.RequestAddFriendStatus `json:"status" bson:"status"`
	// 申请时间
	CreateTime int64 `json:"createTime" bson:"createTime"`
	// 更新时间
	UpdateTime int64 `json:"updateTime" bson:"updateTime"`
	// 附加信息
	Extra []*pb.RequestAddFriendExtra `json:"extra" bson:"extra"`
}

func (m *RequestAddFriend) CollectionName() string {
	return "request_add_friend"
}

func (m *RequestAddFriend) Indexes(c *qmgo.Collection) error {
	_ = c.CreateIndexes(context.Background(), []options.IndexModel{{
		Key: []string{"fromUserId"},
	}, {
		Key: []string{"toUserId"},
	}, {
		Key: []string{"fromUserId", "toUserId"},
	}, {
		Key: []string{"status"},
	}, {
		Key: []string{"updateTime"},
	}})
	return nil
}
