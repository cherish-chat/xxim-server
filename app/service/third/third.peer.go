package third

import (
	"fmt"
	"github.com/cherish-chat/xxim-server/config"

	"github.com/cherish-chat/xxim-proto/peerpb"
	captchaserviceServer "github.com/cherish-chat/xxim-server/app/service/third/internal/server/captchaservice"
	emailserviceServer "github.com/cherish-chat/xxim-server/app/service/third/internal/server/emailservice"
	smsserviceServer "github.com/cherish-chat/xxim-server/app/service/third/internal/server/smsservice"
	"github.com/cherish-chat/xxim-server/app/service/third/internal/svc"

	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Run(cfg config.Config) {
	ctx := svc.NewServiceContext(cfg)

	c := cfg.GetThirdConfig()

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		peerpb.RegisterSmsServiceServer(grpcServer, smsserviceServer.NewSmsServiceServer(ctx))
		peerpb.RegisterEmailServiceServer(grpcServer, emailserviceServer.NewEmailServiceServer(ctx))
		peerpb.RegisterCaptchaServiceServer(grpcServer, captchaserviceServer.NewCaptchaServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
