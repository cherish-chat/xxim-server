package xredis

import (
	"math/rand"
	"time"
)

func ExpireMinutes(min int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(10) + min*60
}

const (
	NotFound = "@@NOT_FOUND@@"
)
