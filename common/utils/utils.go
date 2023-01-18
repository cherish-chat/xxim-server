package utils

import (
	"context"
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
