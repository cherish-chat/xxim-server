package utils

import (
	"encoding/json"
	"google.golang.org/protobuf/proto"
	"strconv"
)

func AnyToString(t any) string {
	switch v := t.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case int32:
		return strconv.FormatInt(int64(v), 10)
	case int16:
		return strconv.FormatInt(int64(v), 10)
	case int8:
		return strconv.FormatInt(int64(v), 10)
	case int:
		return strconv.FormatInt(int64(v), 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case uint32:
		return strconv.FormatUint(uint64(v), 10)
	case uint16:
		return strconv.FormatUint(uint64(v), 10)
	case uint8:
		return strconv.FormatUint(uint64(v), 10)
	case uint:
		return strconv.FormatUint(uint64(v), 10)
	default:
		buf, _ := json.Marshal(v)
		return string(buf)
	}
}

func AnyToBytes(t any) []byte {
	if t == nil {
		return nil
	}
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
