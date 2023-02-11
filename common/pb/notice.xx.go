package pb

import "strings"

func HiddenConvIdCommand() string {
	return NoticeConvId("command")
}

func HiddenConvIdFriendMember() string {
	return NoticeConvId("friendMember")
}

func HiddenConvIdGroupMember() string {
	return NoticeConvId("groupMember")
}

func HiddenConvIdFriend(userId string) string {
	return NoticeConvId("friend@" + userId)
}

func HiddenConvIdGroup(groupId string) string {
	return NoticeConvId(GroupConvId(groupId))
}

func HiddenConvId(convId string) string {
	return NoticeConvId(convId)
}

func HiddenConvIdSingle(selfId string, userId string) string {
	return NoticeConvId(SingleConvId(selfId, userId))
}

// DefaultAckId 默认的ackId -1表示会话最大的ackId
func DefaultAckId(convId string) int64 {
	if convId == HiddenConvIdCommand() || convId == HiddenConvIdFriendMember() || convId == HiddenConvIdGroupMember() {
		return -1
	}
	// 如果包含@，则表示是隐藏会话 也返回-1
	if strings.Contains(convId, "@") {
		return -1
	}
	return 0
}

const (
	// command
	NoticeContentType_SyncFriendList = 101 // 同步好友列表

	// friend@
	NoticeContentType_UpdateUserInfo = 201 // 更新用户信息

	// group@
	NoticeContentType_GroupMemberLeave   = 301 // 群成员离开
	NoticeContentType_CreateGroup        = 302 // 创建群
	NoticeContentType_NewGroupMember     = 303 // 新群成员
	NoticeContentType_DismissGroup       = 304 // 解散群
	NoticeContentType_SetGroupMemberInfo = 305 // 设置群成员信息
	NoticeContentType_SetGroupInfo       = 306 // 设置群信息
	NoticeContentType_RecoverGroup       = 307 // 恢复群

	// groupMember
	NoticeContentType_ApplyToBeGroupMember = 401 // 申请加入群
)

type (
	NoticeContent_UpdateUserInfo struct {
		UserId    string                 `json:"userId"`
		UpdateMap map[string]interface{} `json:"updateMap"`
	}
	NoticeContent_SyncFriendList struct {
		Comment string `json:"comment"`
	}
	NoticeContent_ApplyToBeGroupMember struct {
		ApplyId      string                 `json:"applyId"`
		GroupId      string                 `json:"groupId"`
		UserId       string                 `json:"userId"`
		Result       GroupApplyHandleResult `json:"result"`
		Reason       string                 `json:"reason"`
		ApplyTime    int64                  `json:"applyTime"`
		HandleTime   int64                  `json:"handleTime"`
		HandleUserId string                 `json:"handleUserId"`
	}
	NoticeContent_CreateGroup struct {
		GroupId string `json:"groupId"`
	}
	NoticeContent_NewGroupMember struct {
		GroupId  string `json:"groupId"`
		MemberId string `json:"memberId"`
	}
	NoticeContent_GroupMemberLeave struct {
		GroupId string `json:"groupId"`
		Tip     string `json:"tip"`
	}
	NoticeContent_DismissGroup struct {
		GroupId string `json:"groupId"`
	}
	NoticeContent_SetGroupMemberInfo struct {
		GroupId   string                 `json:"groupId"`
		MemberId  string                 `json:"memberId"`
		UpdateMap map[string]interface{} `json:"updateMap"`
	}
)
