package xredis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"testing"
)

func getrc() redis.UniversalClient {
	return GetClient(Config{
		Host: "localhost:6379",
		Pass: "123456",
		DB:   0,
	})
}
func TestHMSetEx(t *testing.T) {
	HMSetEx(getrc(), context.Background(), "test", map[string]interface{}{
		"1": 1,
		"2": 2,
	}, 600)
}
