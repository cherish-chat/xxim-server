package rediskey

const (
	FriendListExpire    = 60 * 60 * 24
	BlacklistListExpire = 60 * 60 * 24 * 2
)

func FriendList(userId string) string {
	return "h:model:friend:" + userId
}

func BlacklistList(userId string) string {
	return "h:model:blacklist:" + userId
}
