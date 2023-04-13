package main

import (
	"flag"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtservice"
	"github.com/cherish-chat/xxim-server/app/xos/internal/config"
	"github.com/cherish-chat/xxim-server/app/xos/internal/server"
	"github.com/cherish-chat/xxim-server/app/xos/internal/svc"
)

var mgmtRpcAddress = flag.String("a", "127.0.0.1:6708", "mgmt rpc address")

func main() {
	flag.Parse()

	var c config.Config
	mgmtService := mgmtservice.MustLoadConfig(*mgmtRpcAddress, "xos", &c)
	ctx := svc.NewServiceContext(c, mgmtService)
	svr := server.NewHttpServer(ctx)

	svr.Start()
}
