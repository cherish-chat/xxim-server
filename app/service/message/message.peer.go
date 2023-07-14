package message

import (
	"fmt"
	"github.com/cherish-chat/xxim-server/config"

	"github.com/cherish-chat/xxim-proto/peerpb"
	consumeserviceServer "github.com/cherish-chat/xxim-server/app/service/message/internal/server/consumeservice"
	messageserviceServer "github.com/cherish-chat/xxim-server/app/service/message/internal/server/messageservice"
	noticeserviceServer "github.com/cherish-chat/xxim-server/app/service/message/internal/server/noticeservice"
	"github.com/cherish-chat/xxim-server/app/service/message/internal/svc"

	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Run(cfg config.Config) {
	ctx := svc.NewServiceContext(cfg)

	c := cfg.GetMessageConfig()

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		peerpb.RegisterMessageServiceServer(grpcServer, messageserviceServer.NewMessageServiceServer(ctx))
		peerpb.RegisterNoticeServiceServer(grpcServer, noticeserviceServer.NewNoticeServiceServer(ctx))

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	go consumeserviceServer.NewConsumerServer(ctx).Start()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
