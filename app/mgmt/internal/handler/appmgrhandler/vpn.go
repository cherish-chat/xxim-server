package appmgrhandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
)

// getAllVpnList 获取全部appVPN列表
// @Summary 获取全部appVPN列表
// @Description 使用此接口获取全部appVPN列表
// @Tags appVPN管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetAllAppMgmtVpnReq true "请求参数"
// @Success 200 {object} pb.GetAllAppMgmtVpnResp "响应数据"
// @Router /ms/get/vpn/list/all [post]
func (r *AppMgrHandler) getAllVpnList(ctx *gin.Context) {
	in := &pb.GetAllAppMgmtVpnReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().GetAllAppMgmtVpn(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// getVpnDetail 获取appVPN详情
// @Summary 获取appVPN详情
// @Description 使用此接口获取appVPN详情
// @Tags appVPN管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetAppMgmtVpnDetailReq true "请求参数"
// @Success 200 {object} pb.GetAppMgmtVpnDetailResp "响应数据"
// @Router /ms/get/vpn/detail [post]
func (r *AppMgrHandler) getVpnDetail(ctx *gin.Context) {
	in := &pb.GetAppMgmtVpnDetailReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().GetAppMgmtVpnDetail(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// addVpn 新增appVPN
// @Summary 新增appVPN
// @Description 使用此接口新增appVPN
// @Tags appVPN管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.AddAppMgmtVpnReq true "请求参数"
// @Success 200 {object} pb.AddAppMgmtVpnResp "响应数据"
// @Router /ms/add/vpn [post]
func (r *AppMgrHandler) addVpn(ctx *gin.Context) {
	in := &pb.AddAppMgmtVpnReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().AddAppMgmtVpn(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// updateVpn 更新appVPN
// @Summary 更新appVPN
// @Description 使用此接口更新appVPN
// @Tags appVPN管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.UpdateAppMgmtVpnReq true "请求参数"
// @Success 200 {object} pb.UpdateAppMgmtVpnResp "响应数据"
// @Router /ms/update/vpn [post]
func (r *AppMgrHandler) updateVpn(ctx *gin.Context) {
	in := &pb.UpdateAppMgmtVpnReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().UpdateAppMgmtVpn(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// deleteVpn 删除appVPN
// @Summary 删除appVPN
// @Description 使用此接口删除appVPN
// @Tags appVPN管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.DeleteAppMgmtVpnReq true "请求参数"
// @Success 200 {object} pb.DeleteAppMgmtVpnResp "响应数据"
// @Router /ms/delete/vpn [post]
func (r *AppMgrHandler) deleteVpn(ctx *gin.Context) {
	in := &pb.DeleteAppMgmtVpnReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().DeleteAppMgmtVpn(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}
