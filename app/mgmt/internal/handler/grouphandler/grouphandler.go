package grouphandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/gin-gonic/gin"
)

type GroupHandler struct {
	svcCtx *svc.ServiceContext
}

func NewGroupHandler(svcCtx *svc.ServiceContext) *GroupHandler {
	return &GroupHandler{svcCtx: svcCtx}
}

func (r *GroupHandler) Register(g *gin.RouterGroup) {
	group := g.Group("/groupmgmt") // app管理
	{
		// Model 模型
		// 列表
		group.POST("/get/model/list/all", r.getAllModel)
		// 详情
		group.POST("/get/model/detail", r.getModelDetail)
		// 更新
		group.POST("/update/model", r.updateModel)
		// 删除
		group.POST("/dismiss/model", r.dismissModel)
	}
	{
		// Member 群成员
		// 批量插入僵尸号
		group.POST("/insert/member/zombie", r.insertMemberZombie)
		// 清除僵尸号
		group.POST("/clear/member/zombie", r.clearMemberZombie)
	}
}
