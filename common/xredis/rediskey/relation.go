package rediskey

const (
	FriendListExpire    = 60 * 5
	BlacklistListExpire = 60 * 10
)

func FriendList(userId string) string {
	return "h:model:friend:" + userId
}

func BlacklistList(userId string) string {
	return "h:model:blacklist:" + userId
}

func SingleConvSetting(convId string, userId string) string {
	return "s:model:single_conv_setting:" + convId + ":" + userId
}
