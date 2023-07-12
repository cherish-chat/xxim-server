package xcache

import "fmt"

type xRedisVal struct {
	IncrKeyGroupId    string
	HashKeyConvMaxSeq string
	IncrKeyNoticeSort string
}

var RedisVal = &xRedisVal{
	IncrKeyGroupId:    "group:i:group_id",
	HashKeyConvMaxSeq: "max_seq",
	IncrKeyNoticeSort: "notice:i:sort",
}

func (x *xRedisVal) LockKeyUserUsername(username string) string {
	return fmt.Sprintf("user:l:username:%s", username)
}

func (x *xRedisVal) LockKeyUserPhone(phone string, phoneCode string) string {
	return fmt.Sprintf("user:l:phone:%s:%s", phone, phoneCode)
}

func (x *xRedisVal) LockKeyUserEmail(email string) string {
	return fmt.Sprintf("user:l:email:%s", email)
}

func (x *xRedisVal) HashKeyConvKv(convId string, convTyp int32) string {
	return fmt.Sprintf("message:h:conv_kv:%s_%d", convId, convTyp)
}

func (x *xRedisVal) HashKeyConvMinSeq(userId string) string {
	return fmt.Sprintf("min_seq:%s", userId)
}
