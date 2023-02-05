package appmgrhandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
)

// getAllVersionList 获取全部app版本列表
// @Summary 获取全部app版本列表
// @Description 使用此接口获取全部app版本列表
// @Tags app版本管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetAllAppMgmtVersionReq true "请求参数"
// @Success 200 {object} pb.GetAllAppMgmtVersionResp "响应数据"
// @Router /ms/get/version/list/all [post]
func (r *AppMgrHandler) getAllVersionList(ctx *gin.Context) {
	in := &pb.GetAllAppMgmtVersionReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().GetAllAppMgmtVersion(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// getVersionDetail 获取app版本详情
// @Summary 获取app版本详情
// @Description 使用此接口获取app版本详情
// @Tags app版本管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetAppMgmtVersionDetailReq true "请求参数"
// @Success 200 {object} pb.GetAppMgmtVersionDetailResp "响应数据"
// @Router /ms/get/version/detail [post]
func (r *AppMgrHandler) getVersionDetail(ctx *gin.Context) {
	in := &pb.GetAppMgmtVersionDetailReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().GetAppMgmtVersionDetail(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// addVersion 新增app版本
// @Summary 新增app版本
// @Description 使用此接口新增app版本
// @Tags app版本管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.AddAppMgmtVersionReq true "请求参数"
// @Success 200 {object} pb.AddAppMgmtVersionResp "响应数据"
// @Router /ms/add/version [post]
func (r *AppMgrHandler) addVersion(ctx *gin.Context) {
	in := &pb.AddAppMgmtVersionReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().AddAppMgmtVersion(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// updateVersion 更新app版本
// @Summary 更新app版本
// @Description 使用此接口更新app版本
// @Tags app版本管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.UpdateAppMgmtVersionReq true "请求参数"
// @Success 200 {object} pb.UpdateAppMgmtVersionResp "响应数据"
// @Router /ms/update/version [post]
func (r *AppMgrHandler) updateVersion(ctx *gin.Context) {
	in := &pb.UpdateAppMgmtVersionReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().UpdateAppMgmtVersion(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// deleteVersion 删除app版本
// @Summary 删除app版本
// @Description 使用此接口删除app版本
// @Tags app版本管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.DeleteAppMgmtVersionReq true "请求参数"
// @Success 200 {object} pb.DeleteAppMgmtVersionResp "响应数据"
// @Router /ms/delete/version [post]
func (r *AppMgrHandler) deleteVersion(ctx *gin.Context) {
	in := &pb.DeleteAppMgmtVersionReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().DeleteAppMgmtVersion(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}
