package utils

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/bwmarrin/snowflake"
	"math/rand"
	"time"
)

var (
	c           = make(chan string)
	nd          *snowflake.Node
	NodId       int64
	uniquePodId string
)

func init() {
	rand.Seed(time.Now().UnixNano())
	NodId = int64(rand.Intn(128))
	nd, _ = snowflake.NewNode(NodId)
	uniquePodId = nd.Generate().String()
	go loop()
}

func loop() {
	for {
		c <- nd.Generate().String()
	}
}

// GenId 获取ID
func GenId() string {
	return Md5(uniquePodId + <-c)
}

func Md5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	cipher := h.Sum(nil)
	return hex.EncodeToString(cipher)
}
