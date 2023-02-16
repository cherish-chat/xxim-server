package rediskey

import (
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xredis"
)

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

func GroupMemberKey(groupId string, userId string) string {
	return "s:model:group_member:" + groupId + ":" + userId
}

func GroupMemberExpire() int {
	return xredis.ExpireMinutes(5)
}

func GroupMemberSearchKey(groupId string, whereMap map[string]interface{}, whereOrMap map[string][]interface{}, offset int32, limit int32) string {
	suffix := ""
	suffix += "whereMap=" + utils.AnyToString(whereMap)
	suffix += "&"
	suffix += "whereOrMap=" + utils.AnyToString(whereOrMap)
	suffix += "&"
	suffix += "offset=" + utils.AnyToString(offset)
	suffix += "&"
	suffix += "limit=" + utils.AnyToString(limit)
	return "s:list:group_member:search:" + groupId + ":" + utils.Md5(suffix)
}

func GroupMemberSearchKeyList(groupId string) string {
	return "s:list:group_member:search:keys:" + groupId
}

func GroupMemberSearchKeyExpire() int {
	return xredis.ExpireMinutes(60)
}
