package xcache

import "fmt"

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
