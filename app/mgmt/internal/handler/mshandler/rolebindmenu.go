package mshandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/logic"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
)

// bindRoleMenuBatch 绑定角色菜单
// @Summary 绑定角色菜单
// @Description 使用此接口绑定角色菜单
// @Tags 角色菜单绑定相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param object body pb.BindMSRoleMenuReq true "请求参数"
// @Success 200 {object} pb.BindMSRoleMenuResp "响应数据"
// @Router /ms/bind/role/menu/batch [post]
func (r *MSHandler) bindRoleMenuBatch(ctx *gin.Context) {
	in := &pb.BindMSRoleMenuReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewBindMSRoleMenuLogic(ctx, r.svcCtx).BindMSRoleMenu(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	ctx.JSON(200, out)
}

// unbindRoleMenuBatch 解绑角色菜单
// @Summary 解绑角色菜单
// @Description 使用此接口解绑角色菜单
// @Tags 角色菜单绑定相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param object body pb.UnbindMSRoleMenuReq true "请求参数"
// @Success 200 {object} pb.UnbindMSRoleMenuResp "响应数据"
// @Router /ms/unbind/role/menu/batch [post]
func (r *MSHandler) unbindRoleMenuBatch(ctx *gin.Context) {
	in := &pb.UnbindMSRoleMenuReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewUnbindMSRoleMenuLogic(ctx, r.svcCtx).UnbindMSRoleMenu(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	ctx.JSON(200, out)
}
