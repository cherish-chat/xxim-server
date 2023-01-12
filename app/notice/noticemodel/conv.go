package noticemodel

import "github.com/cherish-chat/xxim-server/common/pb"

var (
	ConvId_ConvListChanged    = pb.NoticeConvId("ConvListChanged")
	ConvId_ProfileChanged     = pb.NoticeConvId("ProfileChanged")
	ConvId_ConvInfoChanged    = pb.NoticeConvId("ConvInfoChanged")
	ConvId_ConvSettingChanged = pb.NoticeConvId("ConvSettingChanged")
	ConvId_ConvMemberChanged  = pb.NoticeConvId("ConvMemberChanged")
	ConvId_SyncFriendList     = pb.NoticeConvId("SyncFriendList")
	ConvId_FriendNotice       = pb.NoticeConvId("FriendNotice") // 好友通知
	ConvId_GroupNotice        = pb.NoticeConvId("GroupNotice")  // 群通知
	ConvId_SystemNotice       = pb.NoticeConvId("SystemNotice") // 系统通知
	ConvId_WorldNotice        = pb.NoticeConvId("WorldNotice")  // 惺球通知
)

var (
	DefaultConvIds = []string{
		ConvId_ConvListChanged,
		ConvId_ProfileChanged,
		ConvId_ConvInfoChanged,
		ConvId_ConvSettingChanged,
		ConvId_ConvMemberChanged,
		ConvId_SyncFriendList,
		ConvId_FriendNotice,
		ConvId_GroupNotice,
		ConvId_SystemNotice,
		ConvId_WorldNotice,
	}
)

type (
	ConvType   int32
	JumpType   int32
	NoticeConv struct {
		ConvId   string   `gorm:"column:convId;primary_key;type:varchar(64);not null" json:"convId"`
		ConvName string   `gorm:"column:convName;type:varchar(64);not null" json:"convName"`
		ConvType ConvType `gorm:"column:convType;type:tinyint(4);not null" json:"convType"`
		// 通知号的头像
		Avatar string `gorm:"column:avatar;type:varchar(255);not null" json:"avatar"`
		// 跳转类型
		JumpType JumpType `gorm:"column:jumpType;type:tinyint(4);not null" json:"jumpType"`
		// 跳转路径
		JumpPath string `gorm:"column:jumpPath;type:varchar(255);not null" json:"jumpPath"`
		// 跳转参数
		JumpParam string `gorm:"column:jumpParam;type:varchar(255);not null" json:"jumpParam"`
	}
)

func (m *NoticeConv) TableName() string {
	return "notice_conv"
}
