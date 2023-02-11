package appmgrhandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
)

// getAllEmojiList 获取全部app表情列表
// @Summary 获取全部app表情列表
// @Description 使用此接口获取全部app表情列表
// @Tags app表情管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetAllAppMgmtEmojiReq true "请求参数"
// @Success 200 {object} pb.GetAllAppMgmtEmojiResp "响应数据"
// @Router /ms/get/emoji/list/all [post]
func (r *AppMgrHandler) getAllEmojiList(ctx *gin.Context) {
	in := &pb.GetAllAppMgmtEmojiReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().GetAllAppMgmtEmoji(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// getEmojiDetail 获取app表情详情
// @Summary 获取app表情详情
// @Description 使用此接口获取app表情详情
// @Tags app表情管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetAppMgmtEmojiDetailReq true "请求参数"
// @Success 200 {object} pb.GetAppMgmtEmojiDetailResp "响应数据"
// @Router /ms/get/emoji/detail [post]
func (r *AppMgrHandler) getEmojiDetail(ctx *gin.Context) {
	in := &pb.GetAppMgmtEmojiDetailReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().GetAppMgmtEmojiDetail(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// addEmoji 新增app表情
// @Summary 新增app表情
// @Description 使用此接口新增app表情
// @Tags app表情管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.AddAppMgmtEmojiReq true "请求参数"
// @Success 200 {object} pb.AddAppMgmtEmojiResp "响应数据"
// @Router /ms/add/emoji [post]
func (r *AppMgrHandler) addEmoji(ctx *gin.Context) {
	in := &pb.AddAppMgmtEmojiReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().AddAppMgmtEmoji(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// updateEmoji 更新app表情
// @Summary 更新app表情
// @Description 使用此接口更新app表情
// @Tags app表情管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.UpdateAppMgmtEmojiReq true "请求参数"
// @Success 200 {object} pb.UpdateAppMgmtEmojiResp "响应数据"
// @Router /ms/update/emoji [post]
func (r *AppMgrHandler) updateEmoji(ctx *gin.Context) {
	in := &pb.UpdateAppMgmtEmojiReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().UpdateAppMgmtEmoji(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// deleteEmoji 删除app表情
// @Summary 删除app表情
// @Description 使用此接口删除app表情
// @Tags app表情管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.DeleteAppMgmtEmojiReq true "请求参数"
// @Success 200 {object} pb.DeleteAppMgmtEmojiResp "响应数据"
// @Router /ms/delete/emoji [post]
func (r *AppMgrHandler) deleteEmoji(ctx *gin.Context) {
	in := &pb.DeleteAppMgmtEmojiReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().DeleteAppMgmtEmoji(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}
