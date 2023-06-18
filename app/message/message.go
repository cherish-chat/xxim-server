package main

import (
	"flag"
	"github.com/cherish-chat/xxim-server/app/message/internal/server"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/cherish-chat/xxim-server/app/message/internal/config"
	messageserviceServer "github.com/cherish-chat/xxim-server/app/message/internal/server/messageservice"
	noticeserviceServer "github.com/cherish-chat/xxim-server/app/message/internal/server/noticeservice"
	"github.com/cherish-chat/xxim-server/app/message/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/message.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)
	consumerServer := server.NewConsumerServer(ctx)

	consumerServer.Start()

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterMessageServiceServer(grpcServer, messageserviceServer.NewMessageServiceServer(ctx))
		pb.RegisterNoticeServiceServer(grpcServer, noticeserviceServer.NewNoticeServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	logx.Infof("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
