package relationmodel

import (
	"github.com/cherish-chat/xxim-server/common/pb"
)

type RequestAddFriend struct {
	Id string `json:"id" bson:"_id" gorm:"column:id;primary_key;type:char(32);not null"`
	// 发起人
	FromUserId string `json:"fromUserId" bson:"fromUserId" gorm:"column:fromUserId;type:char(32);not null;index;index:idx_from_user_id_to_user_id;comment:发起人"`
	// 接收人
	ToUserId string `json:"toUserId" bson:"toUserId" gorm:"column:toUserId;type:char(32);not null;index;index:idx_from_user_id_to_user_id;comment:接收人"`
	// 申请状态
	Status pb.RequestAddFriendStatus `json:"status" bson:"status" gorm:"column:status;type:tinyint(4);not null;default:0;index;comment:申请状态"`
	// 申请时间
	CreateTime int64 `json:"createTime" bson:"createTime" gorm:"column:createTime;type:bigint(20);not null;default:0;comment:申请时间"`
	// 更新时间
	UpdateTime int64 `json:"updateTime" bson:"updateTime" gorm:"column:updateTime;type:bigint(20);not null;default:0;index;comment:更新时间"`
	// 附加信息
	Extra pb.RequestAddFriendExtraList `json:"extra" bson:"extra" gorm:"column:extra;type:json;not null;comment:附加信息"`
}

func (m *RequestAddFriend) TableName() string {
	return "request_add_friend"
}
