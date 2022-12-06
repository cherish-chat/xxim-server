package rediskey

func NoticeConvMembersSubscribed(convId string) string {
	return "s:notice_conv_members_subscribed:" + convId
}

func UserAckRecord(userId string, deviceId string) string {
	return "h:user_ack_record:" + userId + ":" + deviceId
}
