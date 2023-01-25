package noticemodel

import "github.com/cherish-chat/xxim-server/common/pb"

var (
	ConvId_ConvListChanged = pb.NoticeConvId("ConvListChanged")
	ConvId_ProfileChanged  = pb.NoticeConvId("ProfileChanged")
	ConvId_SyncFriendList  = pb.NoticeConvId("SyncFriendList")
	ConvId_FriendNotice    = pb.NoticeConvId("FriendNotice") // 好友通知
	ConvId_GroupNotice     = pb.NoticeConvId("GroupNotice")  // 群通知
	ConvId_SystemNotice    = pb.NoticeConvId("SystemNotice") // 系统通知
	ConvId_WorldNotice     = pb.NoticeConvId("WorldNotice")  // 惺球通知
)

var (
	DefaultConvIds = []string{
		ConvId_ConvListChanged,
		ConvId_ProfileChanged,
		ConvId_SyncFriendList,
		ConvId_FriendNotice,
		ConvId_GroupNotice,
		ConvId_SystemNotice,
		ConvId_WorldNotice,
	}
)
