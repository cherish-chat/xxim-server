package rediskey

func UserToken(uid string) string {
	return "h:user_token:" + uid
}
