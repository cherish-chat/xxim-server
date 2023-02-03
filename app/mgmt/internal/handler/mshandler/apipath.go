package mshandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/logic"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
)

// getAllApiPathList 获取所有服务端ApiPath列表
// @Summary 获取所有服务端ApiPath列表
// @Description 使用此接口获取所有服务端ApiPath列表
// @Tags 管理系统服务端ApiPath相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param object body pb.GetAllMSApiPathListReq true "请求参数"
// @Success 200 {object} pb.GetAllMSApiPathListResp "响应数据"
// @Router /ms/get/apipath/list/all [post]
func (r *MSHandler) getAllApiPathList(ctx *gin.Context) {
	in := &pb.GetAllMSApiPathListReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewGetAllMSApiPathListLogic(ctx, r.svcCtx).GetAllMSApiPathList(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// getApiPathList 获取我的服务端ApiPath列表
// @Summary 获取我的服务端ApiPath列表
// @Description 使用此接口获取我的服务端ApiPath列表
// @Tags 管理系统服务端ApiPath相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param object body pb.GetMyMSApiPathListReq true "请求参数"
// @Success 200 {object} pb.GetMyMSApiPathListResp "响应数据"
// @Router /ms/get/apipath/list [post]
func (r *MSHandler) getApiPathList(ctx *gin.Context) {
	in := &pb.GetMyMSApiPathListReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewGetMyMSApiPathListLogic(ctx, r.svcCtx).GetMyMSApiPathList(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// getApiPathDetail 获取服务端ApiPath详情
// @Summary 获取服务端ApiPath详情
// @Description 使用此接口获取服务端ApiPath详情
// @Tags 管理系统服务端ApiPath相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param object body pb.GetMSApiPathDetailReq true "请求参数"
// @Success 200 {object} pb.GetMSApiPathDetailResp "响应数据"
// @Router /ms/get/apipath/detail [post]
func (r *MSHandler) getApiPathDetail(ctx *gin.Context) {
	in := &pb.GetMSApiPathDetailReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewGetMSApiPathDetailLogic(ctx, r.svcCtx).GetMSApiPathDetail(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// addApiPath 新增服务端ApiPath
// @Summary 新增服务端ApiPath
// @Description 使用此接口新增服务端ApiPath
// @Tags 管理系统服务端ApiPath相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param object body pb.AddMSApiPathReq true "请求参数"
// @Success 200 {object} pb.AddMSApiPathResp "响应数据"
// @Router /ms/add/apipath [post]
func (r *MSHandler) addApiPath(ctx *gin.Context) {
	in := &pb.AddMSApiPathReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewAddMSApiPathLogic(ctx, r.svcCtx).AddMSApiPath(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// updateApiPath 更新服务端ApiPath
// @Summary 更新服务端ApiPath
// @Description 使用此接口更新服务端ApiPath
// @Tags 管理系统服务端ApiPath相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param object body pb.UpdateMSApiPathReq true "请求参数"
// @Success 200 {object} pb.UpdateMSApiPathResp "响应数据"
// @Router /ms/update/apipath [post]
func (r *MSHandler) updateApiPath(ctx *gin.Context) {
	in := &pb.UpdateMSApiPathReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewUpdateMSApiPathLogic(ctx, r.svcCtx).UpdateMSApiPath(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// deleteApiPathBatch 删除服务端ApiPath
// @Summary 删除服务端ApiPath
// @Description 使用此接口删除服务端ApiPath
// @Tags 管理系统服务端ApiPath相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param object body pb.DeleteMSApiPathReq true "请求参数"
// @Success 200 {object} pb.DeleteMSApiPathResp "响应数据"
// @Router /ms/delete/apipath [post]
func (r *MSHandler) deleteApiPathBatch(ctx *gin.Context) {
	in := &pb.DeleteMSApiPathReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewDeleteMSApiPathLogic(ctx, r.svcCtx).DeleteMSApiPath(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}
