package gateway

import (
	"github.com/cherish-chat/xxim-server/config"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/cherish-chat/xxim-proto/peerpb"
	connectionserviceServer "github.com/cherish-chat/xxim-server/app/api/gateway/internal/server/connectionservice"
	interfaceserviceServer "github.com/cherish-chat/xxim-server/app/api/gateway/internal/server/interfaceservice"
	internalserviceServer "github.com/cherish-chat/xxim-server/app/api/gateway/internal/server/internalservice"
	"github.com/cherish-chat/xxim-server/app/api/gateway/internal/svc"

	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Run(cfg config.Config) {
	ctx := svc.NewServiceContext(cfg)

	c := cfg.GetGatewayConfig()

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		peerpb.RegisterConnectionServiceServer(grpcServer, connectionserviceServer.NewConnectionServiceServer(ctx))
		peerpb.RegisterInternalServiceServer(grpcServer, internalserviceServer.NewInternalServiceServer(ctx))
		peerpb.RegisterInterfaceServiceServer(grpcServer, interfaceserviceServer.NewInterfaceServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	go interfaceserviceServer.NewCustomInterfaceServiceServer(ctx).Start()

	logx.Infof("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
