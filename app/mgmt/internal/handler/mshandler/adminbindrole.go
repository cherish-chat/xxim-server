package mshandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/logic"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
)

// bindAdminRoleBatch 绑定管理员角色
// @Summary 绑定管理员角色
// @Description 使用此接口绑定管理员角色
// @Tags 管理员角色相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param object body pb.BindMSUserRoleReq true "请求参数"
// @Success 200 {object} pb.BindMSUserRoleResp "响应数据"
// @Router /ms/bind/admin/role/batch [post]
func (r *MSHandler) bindAdminRoleBatch(ctx *gin.Context) {
	in := &pb.BindMSUserRoleReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewBindMSUserRoleLogic(ctx, r.svcCtx).BindMSUserRole(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	ctx.JSON(200, out)
}

// unbindAdminRoleBatch 解绑管理员角色
// @Summary 解绑管理员角色
// @Description 使用此接口解绑管理员角色
// @Tags 管理员角色相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param object body pb.UnbindMSUserRoleReq true "请求参数"
// @Success 200 {object} pb.UnbindMSUserRoleResp "响应数据"
// @Router /ms/unbind/admin/role/batch [post]
func (r *MSHandler) unbindAdminRoleBatch(ctx *gin.Context) {
	in := &pb.UnbindMSUserRoleReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewUnbindMSUserRoleLogic(ctx, r.svcCtx).UnbindMSUserRole(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	ctx.JSON(200, out)
}
