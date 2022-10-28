package main

import (
	"flag"
	"github.com/cherish-chat/xxim-server/app/conn/internal/config"
	"github.com/cherish-chat/xxim-server/app/conn/internal/server"
	"github.com/cherish-chat/xxim-server/app/conn/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "etc/im.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	ctx := svc.NewServiceContext(c)
	svr := server.NewConnServer(ctx)

	svr.Start()
}
