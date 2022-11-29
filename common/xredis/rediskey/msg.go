package rediskey

import "strings"

func ConvKv(convId string) string {
	return "h:conv_kv:" + convId
}

func HKConvMinSeq() string {
	return "min_seq"
}

func HKConvMaxSeq() string {
	return "max_seq"
}

func ConvMsgIdMapping(convId string) string {
	return "h:conv_msgid_mapping:" + convId
}

func MsgKey(id string) string {
	return "s:model:msg:" + id
}

func ConvMembersSubscribed(convId string) string {
	return "s:conv_members_subscribed:" + convId
}

func ConvMemberPodIp(userId string, podIp string) string {
	return userId + "@" + podIp
}

func ConvMembersSubscribedSplit(zmember string) (userId string, podIp string) {
	split := strings.Split(zmember, "@")
	if len(split) == 2 {
		return split[0], split[1]
	}
	return "", ""
}

func OfflinePushMsgListKey(convId string, uid string) string {
	return "s:offline_push_msg:" + convId + ":" + uid
}
