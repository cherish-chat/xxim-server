package pb

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

func HiddenConvIdGroup(userId string) string {
	return NoticeConvId("group@" + userId)
}

const (
	NoticeContentType_SyncFriendList       = 101 // 同步好友列表
	NoticeContentType_ApplyToBeGroupMember = 102 // 申请加入群
	NoticeContentType_GroupMemberLeave     = 103 // 群成员离开
	NoticeContentType_SetGroupMemberInfo   = 104 // 设置群成员信息
	NoticeContentType_UpdateUserInfo       = 201 // 更新用户信息
	NoticeContentType_CreateGroup          = 202 // 创建群
	NoticeContentType_NewGroupMember       = 203 // 新群成员
	NoticeContentType_DismissGroup         = 204 // 解散群
)

type (
	NoticeContent_UpdateUserInfo struct {
		UserId    string                 `json:"userId"`
		UpdateMap map[string]interface{} `json:"updateMap"`
	}
	NoticeContent_SyncFriendList struct {
		Comment string `json:"comment"`
	}
	NoticeContent_ApplyToBeGrouoMember struct {
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
