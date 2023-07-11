package utils

import (
	"math"
	"math/rand"
	"time"
)

type xRandom struct {
}

var Random = &xRandom{}

func (x *xRandom) Int(length int) int {
	rand.Seed(time.Now().UnixNano())
	// 生成随机数 长度为length
	// 比如length为6，那么生成的随机数为100000~999999之间的数
	// start = 10 e length-1次方
	start := math.Pow10(length - 1)
	// end = 10 e length次方 - 1
	end := math.Pow10(length) - 1
	return rand.Intn(int(end)-int(start)) + int(start)
}

func (x *xRandom) SliceString(list []string) string {
	rand.Seed(time.Now().UnixNano())
	return list[rand.Intn(len(list))]
}
