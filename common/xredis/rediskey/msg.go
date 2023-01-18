package rediskey

func ConvKv(convId string) string {
	return "h:conv_kv:" + convId
}

func HKConvMinSeq(userId string) string {
	return "min_seq:" + userId
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

func ConvMemberPodIp(userId string) string {
	return userId
}

func ConvMembersSubscribedSplit(zmember string) string {
	return zmember
}

func OfflinePushMsgListKey(uniqueId string) string {
	return "s:offline_push_msg:" + uniqueId
}

func SyncSendMsgLimiter() string {
	return "token_limter:sync_send_msg_limiter"
}
