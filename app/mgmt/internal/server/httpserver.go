package server

import (
	"fmt"
	_ "github.com/cherish-chat/xxim-server/app/mgmt/docs"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler/appmgrhandler"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler/grouphandler"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler/middleware"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler/msghandler"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler/mshandler"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler/serverhandler"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler/userhandler"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/logic"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
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
	engine.Use(middleware.Recovery())
	engine.Use(middleware.Cors(s.svcCtx.Config.Gin.Cors))
	// routes
	engine.GET("/api/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	apiGroup := engine.Group("/api")
	apiGroup.Use(gin.Logger())
	apiGroup.Use(middleware.Log(s.svcCtx.Mysql()))
	apiGroup.Use(middleware.Auth(s.svcCtx.Redis()))
	apiGroup.Use(middleware.Perms(s.svcCtx.Mysql()))
	serverhandler.NewServerHandler(s.svcCtx).Register(apiGroup)
	mshandler.NewMSHandler(s.svcCtx).Register(apiGroup)
	appmgrhandler.NewAppMgrHandler(s.svcCtx).Register(apiGroup)
	userhandler.NewUserHandler(s.svcCtx).Register(apiGroup)
	grouphandler.NewGroupHandler(s.svcCtx).Register(apiGroup)
	msghandler.NewMsgHandler(s.svcCtx).Register(apiGroup)
	// 表情管理 表情组and表情
	// 配置发现导航中的外链 组and外链
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
	logic.NewInitLogic(s.svcCtx).Init()
}
