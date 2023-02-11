package appmgrhandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
)

// getAllEmojiGroupList 获取全部app表情组列表
// @Summary 获取全部app表情组列表
// @Description 使用此接口获取全部app表情组列表
// @Tags app表情组管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetAllAppMgmtEmojiGroupReq true "请求参数"
// @Success 200 {object} pb.GetAllAppMgmtEmojiGroupResp "响应数据"
// @Router /ms/get/emojigroup/list/all [post]
func (r *AppMgrHandler) getAllEmojiGroupList(ctx *gin.Context) {
	in := &pb.GetAllAppMgmtEmojiGroupReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().GetAllAppMgmtEmojiGroup(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// getEmojiGroupDetail 获取app表情组详情
// @Summary 获取app表情组详情
// @Description 使用此接口获取app表情组详情
// @Tags app表情组管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetAppMgmtEmojiGroupDetailReq true "请求参数"
// @Success 200 {object} pb.GetAppMgmtEmojiGroupDetailResp "响应数据"
// @Router /ms/get/emojigroup/detail [post]
func (r *AppMgrHandler) getEmojiGroupDetail(ctx *gin.Context) {
	in := &pb.GetAppMgmtEmojiGroupDetailReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().GetAppMgmtEmojiGroupDetail(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// updateEmojiGroup 更新app表情组
// @Summary 更新app表情组
// @Description 使用此接口更新app表情组
// @Tags app表情组管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.UpdateAppMgmtEmojiGroupReq true "请求参数"
// @Success 200 {object} pb.UpdateAppMgmtEmojiGroupResp "响应数据"
// @Router /ms/update/emojigroup [post]
func (r *AppMgrHandler) updateEmojiGroup(ctx *gin.Context) {
	in := &pb.UpdateAppMgmtEmojiGroupReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().UpdateAppMgmtEmojiGroup(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// deleteEmojiGroup 删除app表情组
// @Summary 删除app表情组
// @Description 使用此接口删除app表情组
// @Tags app表情组管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.DeleteAppMgmtEmojiGroupReq true "请求参数"
// @Success 200 {object} pb.DeleteAppMgmtEmojiGroupResp "响应数据"
// @Router /ms/delete/emojigroup [post]
func (r *AppMgrHandler) deleteEmojiGroup(ctx *gin.Context) {
	in := &pb.DeleteAppMgmtEmojiGroupReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().DeleteAppMgmtEmojiGroup(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}
