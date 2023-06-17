package main

import (
	"flag"
	"github.com/zeromicro/go-zero/core/logx"

	"github.com/cherish-chat/xxim-server/app/third/internal/config"
	captchaserviceServer "github.com/cherish-chat/xxim-server/app/third/internal/server/captchaservice"
	emailserviceServer "github.com/cherish-chat/xxim-server/app/third/internal/server/emailservice"
	smsserviceServer "github.com/cherish-chat/xxim-server/app/third/internal/server/smsservice"
	"github.com/cherish-chat/xxim-server/app/third/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var configFile = flag.String("f", "etc/third.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterSmsServiceServer(grpcServer, smsserviceServer.NewSmsServiceServer(ctx))
		pb.RegisterEmailServiceServer(grpcServer, emailserviceServer.NewEmailServiceServer(ctx))
		pb.RegisterCaptchaServiceServer(grpcServer, captchaserviceServer.NewCaptchaServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	logx.Infof("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
