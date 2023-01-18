package groupmodel

import (
	"database/sql/driver"
	"encoding/json"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/xorm"
)

type (
	Group struct {
		// 群id 从10001开始 redis incr
		Id     string `bson:"_id" json:"id" gorm:"column:id;primary_key;type:char(32);not null"`
		Name   string `bson:"name" json:"name" gorm:"column:name;type:varchar(32);not null;index;default:''"`
		Avatar string `bson:"avatar" json:"avatar"` // 群头像
		// 群主
		Owner string `bson:"owner" json:"owner" gorm:"column:owner;type:char(32);not null;index"`
		// 所有管理员
		Managers xorm.SliceString `bson:"managers" json:"managers" gorm:"column:managers;type:json;"`
		// 创建时间
		CreateTime int64 `bson:"createTime" json:"createTime" gorm:"column:createTime;type:bigint;not null;index"`
		// 解散时间
		DismissTime int64 `bson:"dismissTime" json:"dismissTime" gorm:"column:dismissTime;type:bigint;not null;index"`
		// 群描述
		Description string `bson:"description" json:"description" gorm:"column:description;type:varchar(255);not null;default:''"`
		// GroupSetting
		Setting GroupSetting `bson:"setting" json:"setting" gorm:"column:setting;type:json;not null"`
	}
	GroupSetting struct {
		// 全体禁言开关
		AllMute bool `bson:"allMute" json:"allMute"`
		// 发言频率限制
		SpeakLimit *int32 `bson:"speakLimit,omitempty" json:"speakLimit"`
		// 群成员人数上限
		MaxMember int32 `bson:"maxMember,omitempty" json:"maxMember"`
		// 成员权限选项
		// 群成员是否可以发起临时会话
		MemberCanStartTempChat bool `bson:"memberCanStartTempChat" json:"memberCanStartTempChat"`
		// 群成员是否可以邀请好友加入群
		MemberCanInviteFriend bool `bson:"memberCanInviteFriend" json:"memberCanInviteFriend"`
		// 新成员可见的历史消息条数
		NewMemberHistoryMsgCount int32 `bson:"newMemberHistoryMsgCount,omitempty" json:"newMemberHistoryMsgCount"`
		// 是否开启匿名聊天
		AnonymousChat   bool            `bson:"anonymousChat" json:"anonymousChat"`
		JoinGroupOption JoinGroupOption `bson:"joinGroupOption" json:"joinGroupOption"`
	}
	JoinGroupOption struct {
		Type pb.GroupSetting_JoinGroupOpt_Type `bson:"type" json:"type"`
		// 验证信息
		// 问题
		Question string `bson:"question" json:"question"`
		// 答案
		Answer string `bson:"answer" json:"answer"`
	}
)

func (m *Group) TableName() string {
	return "group"
}

func (m *Group) GroupBaseInfo() *pb.GroupBaseInfo {
	return &pb.GroupBaseInfo{
		Id:     m.Id,
		Name:   m.Name,
		Avatar: m.Avatar,
	}
}

func (m *Group) Bytes() []byte {
	data, _ := json.Marshal(m)
	return data
}

func (m GroupSetting) Value() (driver.Value, error) {
	return json.Marshal(m)
}

func (m *GroupSetting) Scan(src interface{}) error {
	return json.Unmarshal(src.([]byte), m)
}

func GroupFromBytes(data []byte) *Group {
	group := &Group{}
	err := json.Unmarshal(data, group)
	if err != nil {
		return nil
	}
	return group
}
