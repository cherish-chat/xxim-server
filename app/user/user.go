package main

import (
	"flag"
	"github.com/cherish-chat/xxim-server/app/user/internal/server"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/cherish-chat/xxim-server/app/user/internal/config"
	accountserviceServer "github.com/cherish-chat/xxim-server/app/user/internal/server/accountservice"
	callbackserviceServer "github.com/cherish-chat/xxim-server/app/user/internal/server/callbackservice"
	infoserviceServer "github.com/cherish-chat/xxim-server/app/user/internal/server/infoservice"
	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	consumerServer := server.NewConsumerServer(ctx)

	consumerServer.Start()

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterAccountServiceServer(grpcServer, accountserviceServer.NewAccountServiceServer(ctx))
		pb.RegisterInfoServiceServer(grpcServer, infoserviceServer.NewInfoServiceServer(ctx))
		pb.RegisterCallbackServiceServer(grpcServer, callbackserviceServer.NewCallbackServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	logx.Infof("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
