package main

import (
	"flag"
	"github.com/cherish-chat/xxim-server/app/gateway"
)

var configFile = flag.String("f", "etc/gateway.yaml", "the config file")

func main() {
	flag.Parse()

	gateway.Run(*configFile)
}
