package dbmodel

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/qiniu/qmgo"
	qopt "github.com/qiniu/qmgo/options"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	Msg struct {
		Id          string                  `bson:"_id"`          // server id
		ClientMsgId string                  `bson:"clientMsgId"`  // client id
		ConvId      string                  `bson:"convId"`       // 会话id
		ConvInfo    []byte                  `bson:"convInfo"`     // 会话信息
		SenderId    string                  `bson:"senderId"`     // 发送者id
		SenderInfo  []byte                  `bson:"senderInfo"`   // 发送者信息
		ClientTime  int64                   `bson:"clientTime"`   // 客户端时间
		ServerTime  int64                   `bson:"serverTime"`   // 服务端时间
		Seq         uint32                  `bson:"seq"`          // 消息序号
		ContentType pb.MsgData_ContentType  `bson:"contentType"`  // 消息类型
		OfflinePush *pb.MsgData_OfflinePush `bson:"offline_push"` // 离线推送
		MsgOptions  *pb.MsgData_MsgOptions  `bson:"msg_options"`  // 消息选项
		Ex          []byte                  `bson:"ex,omitempty"` // 扩展字段
		ExcludeUIds []string                `bson:"excludeUIds"`  // 排除的用户id
		DeletedAt   int64                   `bson:"deletedAt"`    // 删除时间
	}
)

func InitMsg(c *qmgo.Collection) {
	c.CreateIndexes(context.Background(), []qopt.IndexModel{{
		Key:          []string{"+convId", "-seq"},
		IndexOptions: options.Index().SetUnique(true),
	}, {
		Key: []string{"+clientMsgId"},
	}})
}
