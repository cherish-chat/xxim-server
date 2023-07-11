package utils

import (
	"github.com/bwmarrin/snowflake"
	"math/rand"
	"time"
)

type xSnowflake struct {
	int64Chan   chan int64
	stringChan  chan string
	nd          *snowflake.Node
	nodId       int64
	uniquePodId string
}

var Snowflake *xSnowflake

func init() {
	rand.Seed(time.Now().UnixNano())
	Snowflake = &xSnowflake{
		int64Chan:  make(chan int64),
		stringChan: make(chan string),
	}
	Snowflake.nodId = int64(rand.Intn(128))
	Snowflake.nd, _ = snowflake.NewNode(Snowflake.nodId)
	Snowflake.uniquePodId = Snowflake.nd.Generate().String()
	Snowflake.loop()
}

func (x *xSnowflake) Int64() int64 {
	return <-x.int64Chan
}

func (x *xSnowflake) String() string {
	return Md5(<-x.stringChan)
}

func (x *xSnowflake) loop() {
	go func() {
		for {
			x.int64Chan <- x.nd.Generate().Int64()
		}
	}()
	go func() {
		for {
			x.stringChan <- x.nd.Generate().String()
		}
	}()
}
