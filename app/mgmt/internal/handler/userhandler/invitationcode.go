package userhandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
)

// getAllInvitationCode 获取全部用户邀请码列表
// @Summary 获取全部用户邀请码列表
// @Description 使用此接口获取全部用户邀请码列表
// @Tags 用户邀请码管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetAllUserInvitationCodeReq true "请求参数"
// @Success 200 {object} pb.GetAllUserInvitationCodeResp "响应数据"
// @Router /ms/get/invitationcode/list/all [post]
func (r *UserHandler) getAllInvitationCode(ctx *gin.Context) {
	in := &pb.GetAllUserInvitationCodeReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.UserService().GetAllUserInvitationCode(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// getInvitationCodeDetail 获取用户邀请码详情
// @Summary 获取用户邀请码详情
// @Description 使用此接口获取用户邀请码详情
// @Tags 用户邀请码管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetUserInvitationCodeDetailReq true "请求参数"
// @Success 200 {object} pb.GetUserInvitationCodeDetailResp "响应数据"
// @Router /ms/get/invitationcode/detail [post]
func (r *UserHandler) getInvitationCodeDetail(ctx *gin.Context) {
	in := &pb.GetUserInvitationCodeDetailReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.UserService().GetUserInvitationCodeDetail(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// addInvitationCode 新增用户邀请码
// @Summary 新增用户邀请码
// @Description 使用此接口新增用户邀请码
// @Tags 用户邀请码管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.AddUserInvitationCodeReq true "请求参数"
// @Success 200 {object} pb.AddUserInvitationCodeResp "响应数据"
// @Router /ms/add/invitationcode [post]
func (r *UserHandler) addInvitationCode(ctx *gin.Context) {
	in := &pb.AddUserInvitationCodeReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.UserService().AddUserInvitationCode(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// updateInvitationCode 更新用户邀请码
// @Summary 更新用户邀请码
// @Description 使用此接口更新用户邀请码
// @Tags 用户邀请码管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.UpdateUserInvitationCodeReq true "请求参数"
// @Success 200 {object} pb.UpdateUserInvitationCodeResp "响应数据"
// @Router /ms/update/invitationcode [post]
func (r *UserHandler) updateInvitationCode(ctx *gin.Context) {
	in := &pb.UpdateUserInvitationCodeReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.UserService().UpdateUserInvitationCode(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// deleteInvitationCode 删除用户邀请码
// @Summary 删除用户邀请码
// @Description 使用此接口删除用户邀请码
// @Tags 用户邀请码管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.DeleteUserInvitationCodeReq true "请求参数"
// @Success 200 {object} pb.DeleteUserInvitationCodeResp "响应数据"
// @Router /ms/delete/invitationcode [post]
func (r *UserHandler) deleteInvitationCode(ctx *gin.Context) {
	in := &pb.DeleteUserInvitationCodeReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.UserService().DeleteUserInvitationCode(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}
