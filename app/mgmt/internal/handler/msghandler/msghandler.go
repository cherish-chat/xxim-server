package msghandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/gin-gonic/gin"
)

type MsgHandler struct {
	svcCtx *svc.ServiceContext
}

func NewMsgHandler(svcCtx *svc.ServiceContext) *MsgHandler {
	return &MsgHandler{svcCtx: svcCtx}
}

func (r *MsgHandler) Register(g *gin.RouterGroup) {
	group := g.Group("/msgmgmt") // app管理
	{
		// Msg 模型
		// 列表
		group.POST("/get/msg/list/all", r.getAllMsg)
	}
}
