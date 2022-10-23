package server

import (
	"github.com/cherish-chat/xxim-server/app/im/internal/svc"
	"github.com/gin-gonic/gin"
)

type ImServiceServer struct {
	svcCtx *svc.ServiceContext
	engine *gin.Engine
}

func NewImServiceServer(svcCtx *svc.ServiceContext) *ImServiceServer {
	return &ImServiceServer{
		svcCtx: svcCtx,
		engine: gin.Default(),
	}
}
