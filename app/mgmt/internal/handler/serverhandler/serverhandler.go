package serverhandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/gin-gonic/gin"
)

type ServerHandler struct {
	svcCtx *svc.ServiceContext
}

func NewServerHandler(svcCtx *svc.ServiceContext) *ServerHandler {
	return &ServerHandler{svcCtx: svcCtx}
}

func (r *ServerHandler) Register(engine *gin.Engine) {
	group := engine.Group("/server")
	// server端配置中心
	{
		group.GET("/config", r.config)
	}
}
