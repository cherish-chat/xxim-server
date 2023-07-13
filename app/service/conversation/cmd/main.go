package main

import (
	"flag"
	"github.com/cherish-chat/xxim-server/app/service/conversation"
	"github.com/cherish-chat/xxim-server/config"
	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "etc/config.yaml", "the config file")

func main() {
	flag.Parse()
	c := config.Config{}
	conf.MustLoad(*configFile, &c)

	conversation.Run(c)
}
