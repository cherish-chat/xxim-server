package userhandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
)

// getAllIpBlackList 获取全部用户ip黑名单列表
// @Summary 获取全部用户ip黑名单列表
// @Description 使用此接口获取全部用户ip黑名单列表
// @Tags 用户ip黑名单管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetAllUserIpBlackListReq true "请求参数"
// @Success 200 {object} pb.GetAllUserIpBlackListResp "响应数据"
// @Router /ms/get/ipblacklist/list/all [post]
func (r *UserHandler) getAllIpBlackList(ctx *gin.Context) {
	in := &pb.GetAllUserIpBlackListReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.UserService().GetAllUserIpBlackList(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// getIpBlackListDetail 获取用户ip黑名单详情
// @Summary 获取用户ip黑名单详情
// @Description 使用此接口获取用户ip黑名单详情
// @Tags 用户ip黑名单管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetUserIpBlackListDetailReq true "请求参数"
// @Success 200 {object} pb.GetUserIpBlackListDetailResp "响应数据"
// @Router /ms/get/ipblacklist/detail [post]
func (r *UserHandler) getIpBlackListDetail(ctx *gin.Context) {
	in := &pb.GetUserIpBlackListDetailReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.UserService().GetUserIpBlackListDetail(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// addIpBlackList 新增用户ip黑名单
// @Summary 新增用户ip黑名单
// @Description 使用此接口新增用户ip黑名单
// @Tags 用户ip黑名单管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.AddUserIpBlackListReq true "请求参数"
// @Success 200 {object} pb.AddUserIpBlackListResp "响应数据"
// @Router /ms/add/ipblacklist [post]
func (r *UserHandler) addIpBlackList(ctx *gin.Context) {
	in := &pb.AddUserIpBlackListReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.UserService().AddUserIpBlackList(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// updateIpBlackList 更新用户ip黑名单
// @Summary 更新用户ip黑名单
// @Description 使用此接口更新用户ip黑名单
// @Tags 用户ip黑名单管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.UpdateUserIpBlackListReq true "请求参数"
// @Success 200 {object} pb.UpdateUserIpBlackListResp "响应数据"
// @Router /ms/update/ipblacklist [post]
func (r *UserHandler) updateIpBlackList(ctx *gin.Context) {
	in := &pb.UpdateUserIpBlackListReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.UserService().UpdateUserIpBlackList(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// deleteIpBlackList 删除用户ip黑名单
// @Summary 删除用户ip黑名单
// @Description 使用此接口删除用户ip黑名单
// @Tags 用户ip黑名单管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.DeleteUserIpBlackListReq true "请求参数"
// @Success 200 {object} pb.DeleteUserIpBlackListResp "响应数据"
// @Router /ms/delete/ipblacklist [post]
func (r *UserHandler) deleteIpBlackList(ctx *gin.Context) {
	in := &pb.DeleteUserIpBlackListReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.UserService().DeleteUserIpBlackList(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}
