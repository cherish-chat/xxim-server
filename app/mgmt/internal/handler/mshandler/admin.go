package mshandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/logic"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
)

// getAllAdminList 获取所有管理员列表
// @Summary 获取所有管理员列表
// @Description 使用此接口获取所有管理员列表
// @Tags 管理员相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param object body pb.GetAllMSUserListReq true "请求参数"
// @Success 200 {object} pb.GetAllMSUserListResp "响应数据"
// @Router /ms/get/admin/list/all [post]
func (r *MSHandler) getAllAdminList(ctx *gin.Context) {
	in := &pb.GetAllMSUserListReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewGetAllMSUserListLogic(ctx, r.svcCtx).GetAllMSUserList(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	ctx.JSON(200, out)
}

// getAdminDetail 获取管理员详情
// @Summary 获取管理员详情
// @Description 使用此接口获取管理员详情
// @Tags 管理员相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param object body pb.GetMSUserDetailReq true "请求参数"
// @Success 200 {object} pb.GetMSUserDetailResp "响应数据"
// @Router /ms/get/admin/detail [post]
func (r *MSHandler) getAdminDetail(ctx *gin.Context) {
	in := &pb.GetMSUserDetailReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewGetMSUserDetailLogic(ctx, r.svcCtx).GetMSUserDetail(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	ctx.JSON(200, out)
}

// addAdmin 添加管理员
// @Summary 添加管理员
// @Description 使用此接口添加管理员
// @Tags 管理员相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param object body pb.AddMSUserReq true "请求参数"
// @Success 200 {object} pb.AddMSUserResp "响应数据"
// @Router /ms/add/admin [post]
func (r *MSHandler) addAdmin(ctx *gin.Context) {
	in := &pb.AddMSUserReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewAddMSUserLogic(ctx, r.svcCtx).AddMSUser(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	ctx.JSON(200, out)
}

// updateAdmin 更新管理员
// @Summary 更新管理员
// @Description 使用此接口更新管理员
// @Tags 管理员相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param object body pb.UpdateMSUserReq true "请求参数"
// @Success 200 {object} pb.UpdateMSUserResp "响应数据"
// @Router /ms/update/admin [post]
func (r *MSHandler) updateAdmin(ctx *gin.Context) {
	in := &pb.UpdateMSUserReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewUpdateMSUserLogic(ctx, r.svcCtx).UpdateMSUser(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	ctx.JSON(200, out)
}

// deleteAdminBatch 删除管理员
// @Summary 删除管理员
// @Description 使用此接口删除管理员
// @Tags 管理员相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param object body pb.DeleteMSUserReq true "请求参数"
// @Success 200 {object} pb.DeleteMSUserResp "响应数据"
// @Router /ms/delete/admin [post]
func (r *MSHandler) deleteAdminBatch(ctx *gin.Context) {
	in := &pb.DeleteMSUserReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewDeleteMSUserLogic(ctx, r.svcCtx).DeleteMSUser(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	ctx.JSON(200, out)
}
