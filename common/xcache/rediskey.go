package xcache

import (
	"fmt"
	"github.com/cherish-chat/xxim-server/common/utils"
)

type xRedisVal struct {
	IncrKeyGroupId           string
	HashKeyConvMessageMaxSeq string
	HashKeyConvNoticeMaxSeq  string
}

var RedisVal = &xRedisVal{
	IncrKeyGroupId:           "group:i:group_id",
	HashKeyConvMessageMaxSeq: "message_max_seq",
	HashKeyConvNoticeMaxSeq:  "notice_max_seq",
}

func (x *xRedisVal) LockKeyUserUsername(username string) string {
	return fmt.Sprintf("user:l:username:%s", username)
}

func (x *xRedisVal) LockKeyUserPhone(phone string, phoneCode string) string {
	return fmt.Sprintf("user:l:phone:%s:%s", phone, phoneCode)
}

func (x *xRedisVal) LockKeyUserEmail(email string) string {
	return fmt.Sprintf("user:l:email:%s:%s", email)
}

func (x *xRedisVal) HashKeyConvKv(convId string, convTyp int32) string {
	return fmt.Sprintf("message:h:conv_kv:%s_%d", convId, convTyp)
}

func (x *xRedisVal) HashKeyConvMessageMinSeq(convId string, userId string, convTyp int32) string {
	// 根据 userId 分 key，避免一个 key 太大
	// 取 utils.Md5(userId) 的后 2 位
	return fmt.Sprintf("message:h:conv_message_min_seq:%s_%s_%d", convId, utils.Md5(userId)[30:], convTyp)
}

func (x *xRedisVal) HashKeyConvNoticeMinSeq(convId string, userId string, convTyp int32) string {
	// 根据 userId 分 key，避免一个 key 太大
	// 取 utils.Md5(userId) 的后 2 位
	return fmt.Sprintf("notice:h:conv_message_min_seq:%s_%s_%d", convId, utils.Md5(userId)[30:], convTyp)
}
