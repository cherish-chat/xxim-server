package main

import (
	"flag"
	"fmt"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/logic"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"log"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/config"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/server"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var redisHost = flag.String("host", "", "redis host")
var redisPass = flag.String("pass", "", "redis password")
var redisTls = flag.Bool("tls", false, "redis tls")
var redisType = flag.String("type", "node", "options: node, cluster")

// go run . -host 127.0.0.1:6803 -pass 123456 -type node

// @title xxim-server HTTP API 文档
// @version 1.0
// @description 此文档由gin-swagger自动生成

// @contact.name showurl

// @host api.cherish.chat
// @BasePath /api
// @schemes https
func main() {
	flag.Parse()

	// 校验配置
	if redisHost == nil || *redisHost == "" {
		log.Fatalf("redis host is empty")
	}

	var c config.Config

	// 连接redis
	redisConf := redis.RedisConf{
		Host: *redisHost,
		Type: *redisType,
		Pass: *redisPass,
		Tls:  *redisTls,
	}
	rc := mgmtmodel.InitRedis(redisConf)
	configFile := logic.MgmtConfig()
	// 打印配置信息
	log.Printf("配置信息: \n%s", configFile)
	if err := conf.LoadFromJsonBytes(configFile, &c); err != nil {
		log.Fatalf("load config failed: %v", err)
	}

	ctx := svc.NewServiceContext(c, rc)
	svr := server.NewMgmtServiceServer(ctx)

	svr.NewHttpServer().Start()

	s := zrpc.MustNewServer(c.RpcServerConf, func(grpcServer *grpc.Server) {
		pb.RegisterMgmtServiceServer(grpcServer, svr)

		if c.Mode == service.DevMode || c.Mode == service.TestMode {
			reflection.Register(grpcServer)
		}
	})
	defer s.Stop()

	fmt.Printf("Starting rpc server at %s...\n", c.ListenOn)
	s.Start()
}
