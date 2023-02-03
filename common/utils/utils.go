package utils

import (
	"context"
	"sync"
	"time"
)

func RetryProxy(ctx context.Context, maxTimes int, interval time.Duration, f func() error) {
	for i := 0; i < maxTimes; i++ {
		select {
		case <-ctx.Done():
			return
		default:
			err := f()
			if err == nil {
				return
			}
			time.Sleep(interval)
		}
	}
}

func SyncMapLength(m *sync.Map) int {
	length := 0
	m.Range(func(key, value interface{}) bool {
		length++
		return true
	})
	return length
}

func TimeFormat(i int64) string {
	// 判断是不是10位
	if i <= 9999999999 && i >= 1000000000 {
		i = i * 1000
	}
	// 判断是不是19位
	if i >= 1000000000000000000 {
		i = i / 1000000
	}
	if i == 0 {
		return ""
	}
	return time.UnixMilli(i).Format("2006-01-02 15:04:05")
}
