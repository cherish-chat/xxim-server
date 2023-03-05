package main

import (
	"flag"
	"fmt"
	"github.com/cherish-chat/xxim-server/app/im/internal/logic"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtservice"

	"github.com/cherish-chat/xxim-server/app/im/internal/config"
	"github.com/cherish-chat/xxim-server/app/im/internal/server"
	"github.com/cherish-chat/xxim-server/app/im/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var mgmtRpcAddress = flag.String("a", "127.0.0.1:6708", "mgmt rpc address")

func main() {
	flag.Parse()

	var c config.Config
	mgmtservice.MustLoadConfig(*mgmtRpcAddress, "im", &c)
	ctx := svc.NewServiceContext(c)
	logic.InitAllIpBlackList(ctx)

	svr := server.NewImServiceServer(ctx)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterImServiceServer(grpcServer, svr)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
