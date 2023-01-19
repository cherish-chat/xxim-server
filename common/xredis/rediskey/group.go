package rediskey

import "github.com/cherish-chat/xxim-server/common/xredis"

func IncrGroup() string {
	return "group"
}

func GroupKey(id string) string {
	return "s:model:group:" + id
}

func GroupKeyExpire() int {
	return xredis.ExpireMinutes(5)
}

func GroupMemberListByUserId(userId string) string {
	return "s:list:group_member:by_user:" + userId
}

func GroupMemberListByUserIdExpire() int {
	return xredis.ExpireMinutes(5)
}
