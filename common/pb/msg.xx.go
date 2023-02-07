package pb

import (
	"strconv"
	"strings"
)

// IdSeparator 会话id之间的分隔符
const IdSeparator = "-"
const SinglePrefix = "single:"
const GroupPrefix = "group:"
const NoticePrefix = "notice:"

func (x *MsgData) IsSingleConv() bool {
	return strings.HasPrefix(x.ConvId, SinglePrefix)
}

func (x *MsgData) ReceiverUid() string {
	split := strings.Split(strings.TrimPrefix(x.ConvId, SinglePrefix), IdSeparator)
	if len(split) == 2 {
		if split[0] == x.SenderId {
			return split[1]
		}
		return split[0]
	}
	return ""
}

func (x *MsgData) ReceiverGid() string {
	return strings.TrimPrefix(x.ConvId, GroupPrefix)
}

func (x *MsgData) IsGroupConv() bool {
	return strings.HasPrefix(x.ConvId, GroupPrefix)
}

func (x *MsgData) IsNoticeConv() bool {
	return strings.HasPrefix(x.ConvId, NoticePrefix)
}

func ServerMsgId(convId string, seq int64) string {
	return convId + IdSeparator + strconv.FormatInt(seq, 10)
}

func ServerNoticeId(convId string, seq int64, userId string) string {
	return convId + IdSeparator + strconv.FormatInt(seq, 10) + IdSeparator + userId
}

func SingleConvId(id1 string, id2 string) string {
	if id1 < id2 {
		return SinglePrefix + id1 + IdSeparator + id2
	}
	return SinglePrefix + id2 + IdSeparator + id1
}

func GroupConvId(groupId string) string {
	return GroupPrefix + groupId
}

func NoticeConvId(noticeId string) string {
	return NoticePrefix + noticeId
}

func IsSingleConv(convId string) bool {
	return strings.HasPrefix(convId, SinglePrefix)
}

func ParseSingleConv(convId string) []string {
	arr := strings.Split(strings.TrimPrefix(convId, SinglePrefix), IdSeparator)
	return arr
}

func ParseServerNoticeId(noticeId string) (convId string, seq int64, userId string) {
	// notice:convId-seq-uid
	noticeId = strings.TrimPrefix(noticeId, NoticePrefix)
	arr := strings.Split(noticeId, IdSeparator)
	// userId 是最后一个
	userId = arr[len(arr)-1]
	// 剩下的是 convId-seq
	convIdSeq := strings.TrimSuffix(noticeId, IdSeparator+userId)
	convId, seq = ParseConvServerMsgId(convIdSeq)
	return
}

func IsGroupConv(convId string) bool {
	return strings.HasPrefix(convId, GroupPrefix)
}

func ParseGroupConv(convId string) string {
	return strings.TrimPrefix(convId, GroupPrefix)
}

func IsNoticeConv(convId string) bool {
	return strings.HasPrefix(convId, NoticePrefix)
}

func ParseConvServerMsgId(serverMsgId string) (convId string, seq int64) {
	arr := strings.Split(serverMsgId, IdSeparator)
	if len(arr) == 2 {
		convId = arr[0]
		seq, _ = strconv.ParseInt(arr[1], 10, 64)
	} else if len(arr) == 3 {
		convId = arr[0] + IdSeparator + arr[1]
		seq, _ = strconv.ParseInt(arr[2], 10, 64)
	}
	return
}
