package groupmodel

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
)

type GroupApply struct {
	Id      string `json:"id" bson:"_id" gorm:"column:id;primary_key;type:varchar(32);not null"`
	GroupId string `json:"groupId" bson:"groupId" gorm:"column:groupId;type:varchar(32);not null;index;"`
	UserId  string `json:"userId" bson:"userId" gorm:"column:userId;type:varchar(32);not null;index;"` // 申请人
	// 申请状态 0:待审核 1:已通过 2:已拒绝
	Result pb.GroupApplyHandleResult `json:"result" bson:"result" gorm:"column:result;type:tinyint(1);not null;default:0;index;"`
	Reason string                    `json:"reason" bson:"reason" gorm:"column:reason;type:varchar(255);not null;default:'';"` // 申请理由
	// 申请时间
	ApplyTime int64 `json:"applyTime" bson:"applyTime" gorm:"column:applyTime;type:bigint(13);not null;default:0;index,sort:desc;"`
	// 处理时间
	HandleTime int64 `json:"handleTime" bson:"handleTime" gorm:"column:handleTime;type:bigint(13);not null;default:0;index;"`
	// 处理人
	HandleUserId string `json:"handleUserId" bson:"handleUserId" gorm:"column:handleUserId;type:varchar(32);not null;default:'';"`
}

func (m *GroupApply) TableName() string {
	return "group_apply"
}

func (m *GroupApply) ToPB() *pb.GroupApplyInfo {
	return &pb.GroupApplyInfo{
		ApplyId:            m.Id,
		GroupId:            m.GroupId,
		UserId:             m.UserId,
		Result:             m.Result,
		Reason:             m.Reason,
		ApplyTime:          m.ApplyTime,
		ApplyTimeStr:       utils.TimeFormat(m.ApplyTime),
		HandleTime:         m.HandleTime,
		HandleTimeStr:      utils.TimeFormat(m.HandleTime),
		HandleUserId:       m.HandleUserId,
		UserBaseInfo:       nil,
		HandleUserBaseInfo: nil,
		GroupBaseInfo:      nil,
	}
}
