package utils

import "strconv"

func AnyMakeSlice[T any](any []T) []T {
	if len(any) == 0 {
		return make([]T, 0)
	}
	return any
}

func AnyPtr[T any](v T) *T {
	return &v
}

func AnyToInt64(t any) int64 {
	switch v := t.(type) {
	case int8:
		return int64(v)
	case int16:
		return int64(v)
	case int32:
		return int64(v)
	case int64:
		return v
	case int:
		return int64(v)
	case uint8:
		return int64(v)
	case uint16:
		return int64(v)
	case uint32:
		return int64(v)
	case uint64:
		return int64(v)
	case uint:
		return int64(v)
	case float32:
		return int64(v)
	case float64:
		return int64(v)
	case string:
		i, _ := strconv.ParseInt(v, 10, 64)
		return i
	default:
		return 0
	}
}
