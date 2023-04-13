package server

import (
	"github.com/cherish-chat/xxim-server/app/xos/internal/handler"
	"github.com/cherish-chat/xxim-server/app/xos/internal/svc"
	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	svcCtx        *svc.ServiceContext
	engine        *gin.Engine
	uploadHandler *handler.UploadHandler
}

func NewHttpServer(svcCtx *svc.ServiceContext) *HttpServer {
	s := &HttpServer{svcCtx: svcCtx}
	s.uploadHandler = handler.NewUploadHandler(svcCtx)
	return s
}

func (s *HttpServer) Start() {
	if s.svcCtx.Config.Mode == "pro" { //dev,pro
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	s.engine = gin.Default()
	s.engine.Use(handler.Cors(s.svcCtx.Config.Gin.Cors))
	s.engine.PUT("/upload/:objectId", s.uploadHandler.PutObject)
	s.engine.POST("/upload/:objectId", s.uploadHandler.PostObject)
	err := s.engine.Run(s.svcCtx.Config.Gin.Addr)
	if err != nil {
		panic(err)
	}
}
