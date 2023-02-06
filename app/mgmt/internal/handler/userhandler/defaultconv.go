package userhandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
)

// getAllDefaultConv 获取全部用户默认会话列表
// @Summary 获取全部用户默认会话列表
// @Description 使用此接口获取全部用户默认会话列表
// @Tags 用户默认会话管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetAllUserDefaultConvReq true "请求参数"
// @Success 200 {object} pb.GetAllUserDefaultConvResp "响应数据"
// @Router /ms/get/defaultconv/list/all [post]
func (r *UserHandler) getAllDefaultConv(ctx *gin.Context) {
	in := &pb.GetAllUserDefaultConvReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.UserService().GetAllUserDefaultConv(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// getDefaultConvDetail 获取用户默认会话详情
// @Summary 获取用户默认会话详情
// @Description 使用此接口获取用户默认会话详情
// @Tags 用户默认会话管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetUserDefaultConvDetailReq true "请求参数"
// @Success 200 {object} pb.GetUserDefaultConvDetailResp "响应数据"
// @Router /ms/get/defaultconv/detail [post]
func (r *UserHandler) getDefaultConvDetail(ctx *gin.Context) {
	in := &pb.GetUserDefaultConvDetailReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.UserService().GetUserDefaultConvDetail(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// addDefaultConv 新增用户默认会话
// @Summary 新增用户默认会话
// @Description 使用此接口新增用户默认会话
// @Tags 用户默认会话管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.AddUserDefaultConvReq true "请求参数"
// @Success 200 {object} pb.AddUserDefaultConvResp "响应数据"
// @Router /ms/add/defaultconv [post]
func (r *UserHandler) addDefaultConv(ctx *gin.Context) {
	in := &pb.AddUserDefaultConvReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.UserService().AddUserDefaultConv(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// updateDefaultConv 更新用户默认会话
// @Summary 更新用户默认会话
// @Description 使用此接口更新用户默认会话
// @Tags 用户默认会话管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.UpdateUserDefaultConvReq true "请求参数"
// @Success 200 {object} pb.UpdateUserDefaultConvResp "响应数据"
// @Router /ms/update/defaultconv [post]
func (r *UserHandler) updateDefaultConv(ctx *gin.Context) {
	in := &pb.UpdateUserDefaultConvReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.UserService().UpdateUserDefaultConv(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// deleteDefaultConv 删除用户默认会话
// @Summary 删除用户默认会话
// @Description 使用此接口删除用户默认会话
// @Tags 用户默认会话管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.DeleteUserDefaultConvReq true "请求参数"
// @Success 200 {object} pb.DeleteUserDefaultConvResp "响应数据"
// @Router /ms/delete/defaultconv [post]
func (r *UserHandler) deleteDefaultConv(ctx *gin.Context) {
	in := &pb.DeleteUserDefaultConvReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.UserService().DeleteUserDefaultConv(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}
