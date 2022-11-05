package rediskey

func ConvKv(convId string) string {
	return "h:conv_kv:" + convId
}

func HKConvMinSeq() string {
	return "min_seq"
}

func HKConvMaxSeq() string {
	return "max_seq"
}

func ConvMsgIdMapping(convId string) string {
	return "h:conv_msgid_mapping:" + convId
}

func MsgKey(id string) string {
	return "s:model:msg:" + id
}
