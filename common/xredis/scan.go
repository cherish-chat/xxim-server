package xredis

import (
	"context"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

func Scan(rc *redis.Redis, ctx context.Context, match string, keys chan []string, stopChan chan struct{}) {
	defer func() {
		close(keys)
		close(stopChan)
	}()
	cursor := uint64(0)
	for {
		var res []string
		var e error
		res, cursor, e = rc.ScanCtx(ctx, cursor, match, 10000)
		keys <- res
		if e != nil {
			return
		}
		if cursor == 0 {
			stopChan <- struct{}{}
			break
		}
	}
}
