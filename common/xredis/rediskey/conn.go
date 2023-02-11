package rediskey

func ConnUserMap(podIp string) string {
	return "h:conn_user_map:" + podIp
}
