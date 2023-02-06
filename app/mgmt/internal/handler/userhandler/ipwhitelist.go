package userhandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
)

// getAllIpWhiteList 获取全部用户ip白名单列表
// @Summary 获取全部用户ip白名单列表
// @Description 使用此接口获取全部用户ip白名单列表
// @Tags 用户ip白名单管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetAllUserIpWhiteListReq true "请求参数"
// @Success 200 {object} pb.GetAllUserIpWhiteListResp "响应数据"
// @Router /ms/get/ipwhitelist/list/all [post]
func (r *UserHandler) getAllIpWhiteList(ctx *gin.Context) {
	in := &pb.GetAllUserIpWhiteListReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.UserService().GetAllUserIpWhiteList(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// getIpWhiteListDetail 获取用户ip白名单详情
// @Summary 获取用户ip白名单详情
// @Description 使用此接口获取用户ip白名单详情
// @Tags 用户ip白名单管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetUserIpWhiteListDetailReq true "请求参数"
// @Success 200 {object} pb.GetUserIpWhiteListDetailResp "响应数据"
// @Router /ms/get/ipwhitelist/detail [post]
func (r *UserHandler) getIpWhiteListDetail(ctx *gin.Context) {
	in := &pb.GetUserIpWhiteListDetailReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.UserService().GetUserIpWhiteListDetail(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// addIpWhiteList 新增用户ip白名单
// @Summary 新增用户ip白名单
// @Description 使用此接口新增用户ip白名单
// @Tags 用户ip白名单管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.AddUserIpWhiteListReq true "请求参数"
// @Success 200 {object} pb.AddUserIpWhiteListResp "响应数据"
// @Router /ms/add/ipwhitelist [post]
func (r *UserHandler) addIpWhiteList(ctx *gin.Context) {
	in := &pb.AddUserIpWhiteListReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.UserService().AddUserIpWhiteList(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// updateIpWhiteList 更新用户ip白名单
// @Summary 更新用户ip白名单
// @Description 使用此接口更新用户ip白名单
// @Tags 用户ip白名单管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.UpdateUserIpWhiteListReq true "请求参数"
// @Success 200 {object} pb.UpdateUserIpWhiteListResp "响应数据"
// @Router /ms/update/ipwhitelist [post]
func (r *UserHandler) updateIpWhiteList(ctx *gin.Context) {
	in := &pb.UpdateUserIpWhiteListReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.UserService().UpdateUserIpWhiteList(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// deleteIpWhiteList 删除用户ip白名单
// @Summary 删除用户ip白名单
// @Description 使用此接口删除用户ip白名单
// @Tags 用户ip白名单管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.DeleteUserIpWhiteListReq true "请求参数"
// @Success 200 {object} pb.DeleteUserIpWhiteListResp "响应数据"
// @Router /ms/delete/ipwhitelist [post]
func (r *UserHandler) deleteIpWhiteList(ctx *gin.Context) {
	in := &pb.DeleteUserIpWhiteListReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.UserService().DeleteUserIpWhiteList(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}
