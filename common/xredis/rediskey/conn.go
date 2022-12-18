package rediskey

func UserConnLock() string {
	return "conn:user_lock"
}

func UserConnLockExpire() int {
	return 10
}
