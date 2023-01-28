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
