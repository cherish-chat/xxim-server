package rediskey

func MQRetryCount(msgId string) string {
	return "s:mq_retry_count:" + msgId
}
