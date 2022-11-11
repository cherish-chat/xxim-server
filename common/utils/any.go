package utils

import (
	"reflect"
	"strconv"
)

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

func Set[T any](slice []T) []T {
	m := make(map[any]struct{})
	for _, v := range slice {
		m[v] = struct{}{}
	}

	var results []T
	for _, v := range slice {
		if _, ok := m[v]; ok {
			results = append(results, v)
			delete(m, v)
		}
	}
	return results
}

func InSlice[T any](slice []T, item T) bool {
	for _, v := range slice {
		if reflect.DeepEqual(v, item) {
			return true
		}
	}
	return false
}

// If 三目运算
func If[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}

func UpdateSlice[T any](slice []T, update func(v T) T) []T {
	var results = make([]T, 0, len(slice))
	for _, v := range slice {
		results = append(results, update(v))
	}
	return results
}

type HasId interface {
	GetId() string
}
type HasIdSerializable interface {
	HasId
	Marshal() []byte
}

func Slice2Map[T HasId](t []T) map[string]T {
	m := make(map[string]T)
	for _, v := range t {
		m[v.GetId()] = v
	}
	return m
}

func Slice2MapBytes[T HasIdSerializable](t []T) map[string][]byte {
	m := make(map[string][]byte)
	for _, v := range t {
		m[v.GetId()] = v.Marshal()
	}
	return m
}

func AnyRandomInSlice[T any](slice []T, defaultValue T) T {
	if len(slice) == 0 {
		return defaultValue
	}
	return slice[RandomInt(0, len(slice))]
}

func SliceRemove[T any](slice []T, value T) []T {
	var newSlice []T
	for _, v := range slice {
		if reflect.DeepEqual(v, value) {
			continue
		}
		newSlice = append(newSlice, v)
	}
	return newSlice
}
