package userhandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/logic"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
)

// getAllModel 获取全部用户模型列表
// @Summary 获取全部用户模型列表
// @Description 使用此接口获取全部用户模型列表
// @Tags 用户模型管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetAllUserModelReq true "请求参数"
// @Success 200 {object} pb.GetAllUserModelResp "响应数据"
// @Router /ms/get/model/list/all [post]
func (r *UserHandler) getAllModel(ctx *gin.Context) {
	in := &pb.GetAllUserModelReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.UserService().GetAllUserModel(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	fileLogic := logic.NewUploadFileLogic(ctx, r.svcCtx)
	for _, model := range out.UserModelList {
		model.Avatar = fileLogic.MayGetUrl(model.Avatar)
	}
	handler.ReturnOk(ctx, out)
}

// getModelDetail 获取用户模型详情
// @Summary 获取用户模型详情
// @Description 使用此接口获取用户模型详情
// @Tags 用户模型管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetUserModelDetailReq true "请求参数"
// @Success 200 {object} pb.GetUserModelDetailResp "响应数据"
// @Router /ms/get/model/detail [post]
func (r *UserHandler) getModelDetail(ctx *gin.Context) {
	in := &pb.GetUserModelDetailReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.UserService().GetUserModelDetail(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// addModel 新增用户模型
// @Summary 新增用户模型
// @Description 使用此接口新增用户模型
// @Tags 用户模型管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.AddUserModelReq true "请求参数"
// @Success 200 {object} pb.AddUserModelResp "响应数据"
// @Router /ms/add/model [post]
func (r *UserHandler) addModel(ctx *gin.Context) {
	in := &pb.AddUserModelReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.UserService().AddUserModel(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// updateModel 更新用户模型
// @Summary 更新用户模型
// @Description 使用此接口更新用户模型
// @Tags 用户模型管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.UpdateUserModelReq true "请求参数"
// @Success 200 {object} pb.UpdateUserModelResp "响应数据"
// @Router /ms/update/model [post]
func (r *UserHandler) updateModel(ctx *gin.Context) {
	in := &pb.UpdateUserModelReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.UserService().UpdateUserModel(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// deleteModel 删除用户模型
// @Summary 删除用户模型
// @Description 使用此接口删除用户模型
// @Tags 用户模型管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.DeleteUserModelReq true "请求参数"
// @Success 200 {object} pb.DeleteUserModelResp "响应数据"
// @Router /ms/delete/model [post]
func (r *UserHandler) deleteModel(ctx *gin.Context) {
	in := &pb.DeleteUserModelReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.UserService().DeleteUserModel(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// switchModel 切换用户状态
// @Summary 切换用户状态
// @Description 使用此接口切换用户状态
// @Tags 用户模型管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.SwitchUserModelReq true "请求参数"
// @Success 200 {object} pb.SwitchUserModelResp "响应数据"
// @Router /ms/switch/model [post]
func (r *UserHandler) switchModel(ctx *gin.Context) {
	in := &pb.SwitchUserModelReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.UserService().SwitchUserModel(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// batchCreateZombie 批量创建僵尸用户
// @Summary 批量创建僵尸用户
// @Description 使用此接口批量创建僵尸用户
// @Tags 用户模型管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.BatchCreateZombieUserReq true "请求参数"
// @Success 200 {object} pb.BatchCreateZombieUserResp "响应数据"
// @Router /ms/batch/create/zombie [post]
func (r *UserHandler) batchCreateZombie(ctx *gin.Context) {
	in := &pb.BatchCreateZombieUserReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.UserService().BatchCreateZombieUser(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}
