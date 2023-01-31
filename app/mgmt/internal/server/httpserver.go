package server

import (
	"fmt"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler/middleware"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler/serverhandler"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/service"
	"log"
)

type HttpServer struct {
	svcCtx *svc.ServiceContext
	*gin.Engine
}

func (s *MgmtServiceServer) NewHttpServer() *HttpServer {
	if s.svcCtx.Config.Mode == service.DevMode {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(middleware.Recovery())
	engine.Use(middleware.Cors(s.svcCtx.Config.Gin.Cors))
	// routes
	serverhandler.NewServerHandler(s.svcCtx).Register(engine)
	return &HttpServer{svcCtx: s.svcCtx, Engine: engine}
}

func (s *HttpServer) Start() {
	go func() {
		fmt.Printf("http server start at %s\n", s.svcCtx.Config.Gin.Addr)
		err := s.Run(s.svcCtx.Config.Gin.Addr)
		if err != nil {
			log.Fatalf("failed to start http server: %v", err)
		}
	}()
}
