package appmgrhandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
)

// getAllLinkList 获取全部app外链列表
// @Summary 获取全部app外链列表
// @Description 使用此接口获取全部app外链列表
// @Tags app外链管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetAllAppMgmtLinkReq true "请求参数"
// @Success 200 {object} pb.GetAllAppMgmtLinkResp "响应数据"
// @Router /ms/get/link/list/all [post]
func (r *AppMgrHandler) getAllLinkList(ctx *gin.Context) {
	in := &pb.GetAllAppMgmtLinkReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().GetAllAppMgmtLink(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// getLinkDetail 获取app外链详情
// @Summary 获取app外链详情
// @Description 使用此接口获取app外链详情
// @Tags app外链管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetAppMgmtLinkDetailReq true "请求参数"
// @Success 200 {object} pb.GetAppMgmtLinkDetailResp "响应数据"
// @Router /ms/get/link/detail [post]
func (r *AppMgrHandler) getLinkDetail(ctx *gin.Context) {
	in := &pb.GetAppMgmtLinkDetailReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().GetAppMgmtLinkDetail(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// addLink 新增app外链
// @Summary 新增app外链
// @Description 使用此接口新增app外链
// @Tags app外链管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.AddAppMgmtLinkReq true "请求参数"
// @Success 200 {object} pb.AddAppMgmtLinkResp "响应数据"
// @Router /ms/add/link [post]
func (r *AppMgrHandler) addLink(ctx *gin.Context) {
	in := &pb.AddAppMgmtLinkReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().AddAppMgmtLink(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// updateLink 更新app外链
// @Summary 更新app外链
// @Description 使用此接口更新app外链
// @Tags app外链管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.UpdateAppMgmtLinkReq true "请求参数"
// @Success 200 {object} pb.UpdateAppMgmtLinkResp "响应数据"
// @Router /ms/update/link [post]
func (r *AppMgrHandler) updateLink(ctx *gin.Context) {
	in := &pb.UpdateAppMgmtLinkReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().UpdateAppMgmtLink(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// deleteLink 删除app外链
// @Summary 删除app外链
// @Description 使用此接口删除app外链
// @Tags app外链管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.DeleteAppMgmtLinkReq true "请求参数"
// @Success 200 {object} pb.DeleteAppMgmtLinkResp "响应数据"
// @Router /ms/delete/link [post]
func (r *AppMgrHandler) deleteLink(ctx *gin.Context) {
	in := &pb.DeleteAppMgmtLinkReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().DeleteAppMgmtLink(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}
