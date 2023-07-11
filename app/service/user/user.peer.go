package user

import (
	"fmt"
	"github.com/cherish-chat/xxim-server/config"

	"github.com/cherish-chat/xxim-proto/peerpb"
	callbackserviceServer "github.com/cherish-chat/xxim-server/app/service/user/internal/server/callbackservice"
	"github.com/cherish-chat/xxim-server/app/service/user/internal/svc"

	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Run(cfg config.Config) {
	ctx := svc.NewServiceContext(cfg)

	c := cfg.GetUserConfig()

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		peerpb.RegisterCallbackServiceServer(grpcServer, callbackserviceServer.NewCallbackServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
