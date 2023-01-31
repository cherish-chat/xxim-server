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

func (r *ServerHandler) Register(g *gin.RouterGroup) {
	group := g.Group("/server")
	// server端配置中心
	{
		group.GET("/config", r.config)
	}
	// 其他
	{
		// 全服在线人数
		group.GET("/onlineshield/:randString", r.onlineShield)
	}
}
