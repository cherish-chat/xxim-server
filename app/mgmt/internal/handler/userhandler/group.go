package userhandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
)

// getAllGroup 群列表
// @Summary 获取群列表
// @Description 使用此接口获取群列表
// @Tags 用户群管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetGroupListByUserIdReq true "请求参数"
// @Success 200 {object} pb.GetGroupListByUserIdResp "响应数据"
// @Router /ms/get/group/list/all [post]
func (r *UserHandler) getAllGroup(ctx *gin.Context) {
	in := &pb.GetGroupListByUserIdReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.GroupService().GetGroupListByUserId(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}
