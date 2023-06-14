package utils

func AnyPtr[T any](v T) *T {
	return &v
}

type EnumInSliceType interface {
	String() string
}
type AnyInSliceType interface {
	string | int8 | int16 | int32 | int64 | int | uint8 | uint16 | uint32 | uint64 | uint | float32 | float64 | bool
}

func EnumInSlice[T EnumInSliceType](v T, slice []T) bool {
	for _, item := range slice {
		if item.String() == v.String() {
			return true
		}
	}
	return false
}

func AnyInSlice[T AnyInSliceType](v T, slice []T) bool {
	for _, item := range slice {
		if item == v {
			return true
		}
	}
	return false
}

func AnyString(o any) string {
	switch v := o.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	case *string:
		if v == nil {
			return ""
		}
		return *v
	case int8:
		return Number.Int64ToString(int64(v))
	case int16:
		return Number.Int64ToString(int64(v))
	case int32:
		return Number.Int64ToString(int64(v))
	case int64:
		return Number.Int64ToString(v)
	case int:
		return Number.Int64ToString(int64(v))
	case uint8:
		return Number.Int64ToString(int64(v))
	case uint16:
		return Number.Int64ToString(int64(v))
	case uint32:
		return Number.Int64ToString(int64(v))
	case uint64:
		return Number.Int64ToString(int64(v))
	case uint:
		return Number.Int64ToString(int64(v))
	case float32:
		return Number.Float64ToString(float64(v))
	case float64:
		return Number.Float64ToString(v)
	case bool:
		if v {
			return "true"
		} else {
			return "false"
		}
	default:
		return Json.MarshalToString(v)
	}
}
