package rediskey

func IncrId() string {
	return "incrId"
}

// LatestConnectRecord : 最近一次连接记录
func LatestConnectRecord(userId string) string {
	return "s:latest_connect_record:" + userId
}

// LatestConnectRecordExpire : 最近一次连接记录过期时间
func LatestConnectRecordExpire() int {
	// 1年
	return 60 * 60 * 24 * 365
}
