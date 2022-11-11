package rediskey

func UserToken(uid string) string {
	return "h:user_token:" + uid
}

func UserKey(id string) string {
	return "s:model:user:" + id
}
