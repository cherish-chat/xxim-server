package utils

import (
	"math/rand"
	"time"
)

func RandomInt(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}
