package rediskey

import "github.com/cherish-chat/xxim-server/common/xredis"

func UserAckRecord(userId string, deviceId string) string {
	return "h:user_ack_record:" + userId + ":" + deviceId
}

func NoticeSortSetKey(convId string, userId string, deviceId string) string {
	return "z:notice:" + convId + ":" + userId + ":" + deviceId
}

func NoticeConvAutoId(convId string) string {
	return "h:notice_conv_auto_id:" + convId
}

func NoticeSortSetExpire() int {
	return xredis.ExpireMinutes(5)
}

func NoticeStringKey(id string) string {
	return "s:model:notice:" + id
}

func NoticeStringExpire() int {
	return xredis.ExpireMinutes(5)
}
