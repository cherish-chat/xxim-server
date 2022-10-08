package jwt

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/xredis"
	"github.com/go-redis/redis/v8"
	"testing"
	"time"
)

var ctx = context.Background()

func TestJwt(t *testing.T) {
	token := GenerateTokenButNotSet("123")
	t.Logf("token: %s", token)
	err := SetToken(ctx, getrc(), token, "Platform-Test")
	if err != nil {
		t.Fatalf("SetToken failed: %v", err)
	}
	code, tip := VerifyToken(ctx, getrc(), token, "Platform-Test")
	if code != VerifyTokenCodeOK {
		t.Fatalf("VerifyToken failed: %v", tip)
	}
	code = BanToken(ctx, getrc(), "123", "你已被禁言", time.Now().Add(time.Second*10))
	if code != VerifyTokenCodeBaned {
		t.Fatalf("BanToken failed: %v", tip)
	}
	code, tip = VerifyToken(ctx, getrc(), token, "Platform-Test")
	if code != VerifyTokenCodeOK {
		t.Logf("VerifyToken failed: %v", tip)
	}
	err = SetToken(ctx, getrc(), token, "Platform-Test")
	if err != nil {
		t.Logf("SetToken failed: %v", err)
	}
	time.Sleep(time.Second * 10)
	code, tip = VerifyToken(ctx, getrc(), token, "Platform-Test")
	if code != VerifyTokenCodeBaned {
		t.Logf("VerifyToken ok: %v", tip)
	}
}

func getrc() redis.UniversalClient {
	return xredis.GetClient(xredis.Config{
		Host: "localhost:6379",
		Pass: "123456",
		DB:   7,
	})
}
