package utils

import (
	"math/rand"
	"strings"
	"time"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type H map[string]any

func Int64Ptr(i int64) *int64 {
	return &i
}

func Conditional[T any](flag bool, a T, b T) T {
	if flag {
		return a
	} else {
		return b
	}
}

func RandString(length int) string {
	var result strings.Builder
	for i := 0; i < length; i++ {
		rand.Seed(time.Now().UnixNano())
		result.WriteByte(chars[rand.Intn(len(chars))])
	}
	return result.String()
}
