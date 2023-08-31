package msgmodel

import (
	"encoding/json"
	"github.com/cherish-chat/xxim-server/common/pb"
)

/*
	message RedPacket {
	  string redPacketId = 1;
	  string convId = 2;
	  string title = 3;
	  //红包类型
	  RedPacketType redPacketType = 4;
	  //红包个数
	  int32 count = 5;
	  //红包总金额 单位分 只有RedPacketType=Random_RedPacket时才有用
	  int64 totalAmount = 6;
	  //单个红包金额 单位分 只有RedPacketType=Normal_RedPacket时才有用
	  int64 singleAmount = 7;
	  //红包封面
	  string cover = 8;
	  //红包领取状态
	  RedPacketStatus redPacketStatus = 9;
	  //红包领取者
	  message Receiver {
	    string userId = 1;
	    int64 amount = 2;
	    int64 receiveTime = 3;
	    string avatar = 4;
	    string nickName = 5;
	  }
	  repeated Receiver receiverList = 10;
	  //红包发送时间
	  int64 sendTime = 11;
	  //红包过期时间
	  int64 expireTime = 12;
	}
*/
type RedPacket struct {
	RedPacketId     string             `gorm:"column:redPacketId;primary_key;type:char(32);not null" json:"redPacketId"`
	SenderId        string             `gorm:"column:senderId;type:char(32);not null" json:"senderId"`
	SenderInfo      string             `gorm:"column:senderInfo;type:varchar(255);not null" json:"senderInfo"`
	ConvId          string             `gorm:"column:convId;type:char(64);not null;index;" json:"convId"`
	Title           string             `gorm:"column:title;type:varchar(255);not null" json:"title"`
	RedPacketType   pb.RedPacketType   `gorm:"column:redPacketType;type:tinyint(4);not null" json:"redPacketType"`
	Count           int32              `gorm:"column:count;type:int(11);not null" json:"count"`
	TotalAmount     int64              `gorm:"column:totalAmount;type:bigint(20);not null" json:"totalAmount"`
	SingleAmount    int64              `gorm:"column:singleAmount;type:bigint(20);not null" json:"singleAmount"`
	Cover           string             `gorm:"column:cover;type:varchar(255);not null" json:"cover"`
	RedPacketStatus pb.RedPacketStatus `gorm:"column:redPacketStatus;type:tinyint(4);not null;index;" json:"redPacketStatus"`
	ReceiverList    string             `gorm:"column:receiverList;type:text;" json:"receiverList"`
	SendTime        int64              `gorm:"column:sendTime;type:bigint(20);not null" json:"sendTime"`
	ExpireTime      int64              `gorm:"column:expireTime;type:bigint(20);not null;index;" json:"expireTime"`

	TransactionId string `gorm:"column:transactionId;type:char(32);not null;index;default:'';" json:"transactionId"`
}

func (m *RedPacket) TableName() string {
	return "red_packet"
}

func (m *RedPacket) GetReceiverList() []*pb.RedPacket_Receiver {
	var list []*pb.RedPacket_Receiver
	_ = json.Unmarshal([]byte(m.ReceiverList), &list)
	return list
}
