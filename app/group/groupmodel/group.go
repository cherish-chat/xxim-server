package groupmodel

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/qiniu/qmgo"
)

type (
	Group struct {
		Id     string `bson:"_id" json:"id"`        // 群id 从10001开始 redis incr
		Name   string `bson:"name" json:"name"`     // 群名称
		Avatar string `bson:"avatar" json:"avatar"` // 群头像
		// 群主
		Owner string `bson:"owner" json:"owner"`
		// 所有管理员
		Managers []string `bson:"managers" json:"managers"`
		// 创建时间
		CreateTime int64 `bson:"createTime" json:"createTime"`
		// 解散时间
		DismissTime int64 `bson:"dismissTime" json:"dismissTime"`
		// 群描述
		Description string `bson:"description" json:"description"`
		// GroupSetting
		Setting GroupSetting `bson:"setting" json:"setting"`
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

func (m *Group) CollectionName() string {
	return "group"
}

func (m *Group) Indexes(c *qmgo.Collection) error {
	// TODO indexes
	return nil
}
