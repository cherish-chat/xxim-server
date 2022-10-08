package rediskey

func AuthToken(midToken string) string {
	return "user:token:" + midToken
}

func LoginLock() string {
	return "user:login:lock"
}
