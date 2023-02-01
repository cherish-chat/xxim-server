package mshandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/logic"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
)

// bindRoleApiPathBatch 绑定角色ApiPath
// @Summary 绑定角色ApiPath
// @Description 使用此接口绑定角色ApiPath
// @Tags 角色ApiPath绑定相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param object body pb.BindMSRoleApiPathReq true "请求参数"
// @Success 200 {object} pb.BindMSRoleApiPathResp "响应数据"
// @Router /ms/bind/role/apipath/batch [post]
func (r *MSHandler) bindRoleApiPathBatch(ctx *gin.Context) {
	in := &pb.BindMSRoleApiPathReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewBindMSRoleApiPathLogic(ctx, r.svcCtx).BindMSRoleApiPath(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	ctx.JSON(200, out)
}

// unbindRoleApiPathBatch 解绑角色ApiPath
// @Summary 解绑角色ApiPath
// @Description 使用此接口解绑角色ApiPath
// @Tags 角色ApiPath绑定相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param object body pb.UnbindMSRoleApiPathReq true "请求参数"
// @Success 200 {object} pb.UnbindMSRoleApiPathResp "响应数据"
// @Router /ms/unbind/role/apipath/batch [post]
func (r *MSHandler) unbindRoleApiPathBatch(ctx *gin.Context) {
	in := &pb.UnbindMSRoleApiPathReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewUnbindMSRoleApiPathLogic(ctx, r.svcCtx).UnbindMSRoleApiPath(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	ctx.JSON(200, out)
}
