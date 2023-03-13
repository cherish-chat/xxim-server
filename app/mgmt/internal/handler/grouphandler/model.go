package grouphandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/logic"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
)

// getAllModel 获取全部群聊模型列表
// @Summary 获取全部群聊模型列表
// @Description 使用此接口获取全部群聊模型列表
// @Tags 群聊模型管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "群聊令牌"
// @Param GroupId header string true "群聊ID"
// @Param object body pb.GetAllGroupModelReq true "请求参数"
// @Success 200 {object} pb.GetAllGroupModelResp "响应数据"
// @Router /ms/get/model/list/all [post]
func (r *GroupHandler) getAllModel(ctx *gin.Context) {
	in := &pb.GetAllGroupModelReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.GroupService().GetAllGroupModel(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	fileLogic := logic.NewUploadFileLogic(ctx, r.svcCtx)
	for _, model := range out.GroupModels {
		model.Avatar = fileLogic.MayGetUrl(model.Avatar)
	}
	handler.ReturnOk(ctx, out)
}

// getModelDetail 获取群聊模型详情
// @Summary 获取群聊模型详情
// @Description 使用此接口获取群聊模型详情
// @Tags 群聊模型管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "群聊令牌"
// @Param GroupId header string true "群聊ID"
// @Param object body pb.GetGroupModelDetailReq true "请求参数"
// @Success 200 {object} pb.GetGroupModelDetailResp "响应数据"
// @Router /ms/get/model/detail [post]
func (r *GroupHandler) getModelDetail(ctx *gin.Context) {
	in := &pb.GetGroupModelDetailReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.GroupService().GetGroupModelDetail(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// updateModel 更新群聊模型
// @Summary 更新群聊模型
// @Description 使用此接口更新群聊模型
// @Tags 群聊模型管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "群聊令牌"
// @Param GroupId header string true "群聊ID"
// @Param object body pb.UpdateGroupModelReq true "请求参数"
// @Success 200 {object} pb.UpdateGroupModelResp "响应数据"
// @Router /ms/update/model [post]
func (r *GroupHandler) updateModel(ctx *gin.Context) {
	in := &pb.UpdateGroupModelReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.GroupService().UpdateGroupModel(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// dismissModel 解散群聊模型
// @Summary 切换群聊状态
// @Description 使用此接口切换群聊状态
// @Tags 群聊模型管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "群聊令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.DismissGroupModelReq true "请求参数"
// @Success 200 {object} pb.DismissGroupModelResp "响应数据"
// @Router /ms/dismiss/model [post]
func (r *GroupHandler) dismissModel(ctx *gin.Context) {
	in := &pb.DismissGroupModelReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.GroupService().DismissGroupModel(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}
