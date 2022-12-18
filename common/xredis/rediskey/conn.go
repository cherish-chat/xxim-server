package rediskey

func UserConnLock(userId string) string {
	return "conn:user_lock:" + userId
}

func UserConnLockExpire() int {
	return 10
}
