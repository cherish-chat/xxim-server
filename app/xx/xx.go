package main

import (
	"flag"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/cherish-chat/xxim-server/app/xx/internal/config"
	"github.com/cherish-chat/xxim-server/app/xx/internal/server"
	"github.com/cherish-chat/xxim-server/app/xx/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/local.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	logx.DisableStat()

	ctx := svc.NewServiceContext(c)
	svr := server.NewXxServiceServer(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterXxServiceServer(grpcServer, svr)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
