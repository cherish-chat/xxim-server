package main

import (
	"flag"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/cherish-chat/xxim-server/app/conversation/internal/config"
	conversationserviceServer "github.com/cherish-chat/xxim-server/app/conversation/internal/server/conversationservice"
	friendserviceServer "github.com/cherish-chat/xxim-server/app/conversation/internal/server/friendservice"
	groupserviceServer "github.com/cherish-chat/xxim-server/app/conversation/internal/server/groupservice"
	"github.com/cherish-chat/xxim-server/app/conversation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/conversation.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterGroupServiceServer(grpcServer, groupserviceServer.NewGroupServiceServer(ctx))
		pb.RegisterFriendServiceServer(grpcServer, friendserviceServer.NewFriendServiceServer(ctx))
		pb.RegisterConversationServiceServer(grpcServer, conversationserviceServer.NewConversationServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	logx.Infof("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
