package friendmodel

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	opts "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type FriendApplyStatus = pb.FriendApplyStatus

const (
	FriendApplyStatusApplying = pb.FriendApplyStatus_Applying
	FriendApplyStatusAccepted = pb.FriendApplyStatus_Accepted
	FriendApplyStatusRejected = pb.FriendApplyStatus_Rejected
)

// FriendApplyRecord 好友申请记录 数据库模型
type FriendApplyRecord struct {
	ApplyId string `bson:"applyId" json:"applyId"` // 申请id
	FromId  string `bson:"fromId" json:"fromId"`   // 申请人id
	ToId    string `bson:"toId" json:"toId"`       // 被申请人id
	//附加验证信息
	Message string `bson:"message" json:"message"` // 附加验证信息
	//附加问题答案
	Answer string `bson:"answer" json:"answer"` // 附加答案
	//申请时间
	ApplyTime primitive.DateTime `bson:"applyTime" json:"applyTime"` // 申请时间
	//申请状态
	Status FriendApplyStatus `bson:"status" json:"status"` // 申请状态
	//From是否将此次请求删除
	FromDeleteTime primitive.DateTime `bson:"fromDeleteTime" json:"fromDeleteTime"` // From是否将此次请求删除
	//To是否将此次请求删除
	ToDeleteTime primitive.DateTime `bson:"toDeleteTime" json:"toDeleteTime"` // To是否将此次请求删除
}

func (m *FriendApplyRecord) GetIndexes() []opts.IndexModel {
	return []opts.IndexModel{{
		Key:          []string{"applyId"},
		IndexOptions: options.Index().SetUnique(true),
	}, {
		Key: []string{"fromId"},
	}, {
		Key: []string{"toId"},
	}, {
		Key: []string{"status"},
	}, {
		Key: []string{"fromDeleteTime"},
	}, {
		Key: []string{"toDeleteTime"},
	}}
}
