package main

import (
	"flag"
	"fmt"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtservice"
	"github.com/cherish-chat/xxim-server/app/msg/internal/logic"

	"github.com/cherish-chat/xxim-server/app/msg/internal/config"
	"github.com/cherish-chat/xxim-server/app/msg/internal/server"
	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"net/http"
	_ "net/http/pprof"
)

var mgmtRpcAddress = flag.String("a", "127.0.0.1:6708", "mgmt rpc address")

func main() {
	go func() {
		http.ListenAndServe(":6060", nil)
	}()
	flag.Parse()

	var c config.Config
	mgmtservice.MustLoadConfig(*mgmtRpcAddress, "msg", &c)
	ctx := svc.NewServiceContext(c)
	logic.InitShieldWordTrieTree(ctx)

	go logic.NewTimerCleaner(ctx).Start()

	svr := server.NewMsgServiceServer(ctx)
	svr.Start()

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterMsgServiceServer(grpcServer, svr)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
