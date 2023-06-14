package utils

import "strconv"

type xNumber struct {
}

var Number = &xNumber{}

func (x *xNumber) Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func (x *xNumber) Float64ToString(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func (x *xNumber) Any2Int64(v any) int64 {
	switch v.(type) {
	case int:
		return int64(v.(int))
	case int8:
		return int64(v.(int8))
	case int16:
		return int64(v.(int16))
	case int32:
		return int64(v.(int32))
	case int64:
		return v.(int64)
	case uint:
		return int64(v.(uint))
	case uint8:
		return int64(v.(uint8))
	case uint16:
		return int64(v.(uint16))
	case uint32:
		return int64(v.(uint32))
	case uint64:
		return int64(v.(uint64))
	case float32:
		return int64(v.(float32))
	case float64:
		return int64(v.(float64))
	case string:
		i, _ := strconv.ParseInt(v.(string), 10, 64)
		return i
	default:
		return 0
	}
}
