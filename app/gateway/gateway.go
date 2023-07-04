package main

import (
	"flag"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/server"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/cherish-chat/xxim-server/app/gateway/internal/config"
	gatewayserviceServer "github.com/cherish-chat/xxim-server/app/gateway/internal/server/gatewayservice"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/gateway.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterGatewayServiceServer(grpcServer, gatewayserviceServer.NewGatewayServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()
	httpServer := server.NewHttpServer(ctx)
	go httpServer.Start()
	rtcServer := server.NewRtcServer(ctx)
	go rtcServer.Start()
	logx.Infof("rpc server start at %s", c.ListenOn)
	s.Start()
}
