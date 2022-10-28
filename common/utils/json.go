package utils

import (
	"encoding/json"
	"google.golang.org/protobuf/proto"
)

func AnyToString(t any) string {
	switch v := t.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	default:
		buf, _ := json.Marshal(v)
		return string(buf)
	}
}

func AnyToBytes(t any) []byte {
	switch v := t.(type) {
	case string:
		return []byte(v)
	case []byte:
		return v
	default:
		buf, _ := json.Marshal(v)
		return buf
	}
}

func ProtoToBytes(t proto.Message) []byte {
	buf, _ := proto.Marshal(t)
	return buf
}
