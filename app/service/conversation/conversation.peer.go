package conversation

import (
	"fmt"
	"github.com/cherish-chat/xxim-proto/peerpb"
	channelserviceServer "github.com/cherish-chat/xxim-server/app/service/conversation/internal/server/channelservice"
	friendserviceServer "github.com/cherish-chat/xxim-server/app/service/conversation/internal/server/friendservice"
	groupserviceServer "github.com/cherish-chat/xxim-server/app/service/conversation/internal/server/groupservice"
	sessionserviceServer "github.com/cherish-chat/xxim-server/app/service/conversation/internal/server/sessionservice"
	"github.com/cherish-chat/xxim-server/app/service/conversation/internal/svc"
	"github.com/cherish-chat/xxim-server/config"

	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Run(cfg config.Config) {
	ctx := svc.NewServiceContext(cfg)

	c := cfg.GetConversationConfig()

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		peerpb.RegisterGroupServiceServer(grpcServer, groupserviceServer.NewGroupServiceServer(ctx))
		peerpb.RegisterFriendServiceServer(grpcServer, friendserviceServer.NewFriendServiceServer(ctx))
		peerpb.RegisterChannelServiceServer(grpcServer, channelserviceServer.NewChannelServiceServer(ctx))
		peerpb.RegisterSessionServiceServer(grpcServer, sessionserviceServer.NewSessionServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
