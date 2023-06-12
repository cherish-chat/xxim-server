package utils

func AnyPtr[T any](v T) *T {
	return &v
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
