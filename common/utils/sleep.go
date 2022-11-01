package utils

import (
	"context"
	"sync"
	"time"
)

var sleepTask sync.Map

func init() {
	// 定时删除所有的 sleepTask
	go func() {
		ticker := time.NewTicker(time.Second * 10)
		for {
			select {
			case <-ticker.C:
				sleepTask.Range(func(key, value interface{}) bool {
					sleepTask.Delete(key)
					return true
				})
			}
		}
	}()
}

func getDelay(times int) time.Duration {
	return time.Duration(times%10) * time.Second
}

func ProxySleep(ctx context.Context) {
	// 判断是否超时
	select {
	case <-ctx.Done():
		return
	default:
		// sleep
		times, ok := sleepTask.Load(ctx)
		if !ok {
			sleepTask.Store(ctx, 1)
			time.Sleep(getDelay(1))
		} else {
			sleepTask.Store(ctx, times.(int)+1)
			time.Sleep(getDelay(times.(int)))
		}
	}
}
