package server

import (
	"fmt"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/handler"
	gatewayservicelogic "github.com/cherish-chat/xxim-server/app/gateway/internal/logic/gatewayservice"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/middleware"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
	"os"
)

type logxWriter struct {
}

func (l *logxWriter) Write(p []byte) (n int, err error) {
	logx.Infof("%s", p)
	return len(p), nil
}

func NewHttpServer(svcCtx *svc.ServiceContext) *HttpServer {
	s := &HttpServer{svcCtx: svcCtx}
	gin.DefaultWriter = new(logxWriter)
	s.ginEngine = gin.New()
	s.ginEngine.Use(
		middleware.Tracing(svcCtx), // 链路追踪
		middleware.Logger(svcCtx),  // 访问日志
		gin.Recovery(),             // panic 恢复
		middleware.Cors(svcCtx),    // 跨域
		middleware.ApiLog(svcCtx),  // api 日志
	)
	handler.SetupRoutes(s.svcCtx, s.ginEngine)
	gatewayservicelogic.InitConnectionLogic(s.svcCtx)
	if s.svcCtx.Config.Mode != "pro" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}
	return s
}

type HttpServer struct {
	svcCtx    *svc.ServiceContext
	ginEngine *gin.Engine
}

func (s *HttpServer) Start() {
	listenOn := fmt.Sprintf("%s:%d", s.svcCtx.Config.Http.Host, s.svcCtx.Config.Http.Port)
	logx.Infof("http server start at %s", listenOn)
	err := s.ginEngine.Run(listenOn)
	if err != nil {
		logx.Errorf("ginEngine.Run error: %v", err)
		os.Exit(1)
	}
}
