package appmgrhandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
)

// getAllRichArticleList 获取全部app富文章列表
// @Summary 获取全部app富文章列表
// @Description 使用此接口获取全部app富文章列表
// @Tags app富文章管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetAllAppMgmtRichArticleReq true "请求参数"
// @Success 200 {object} pb.GetAllAppMgmtRichArticleResp "响应数据"
// @Router /ms/get/richArticle/list/all [post]
func (r *AppMgrHandler) getAllRichArticleList(ctx *gin.Context) {
	in := &pb.GetAllAppMgmtRichArticleReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().GetAllAppMgmtRichArticle(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// getRichArticleDetail 获取app富文章详情
// @Summary 获取app富文章详情
// @Description 使用此接口获取app富文章详情
// @Tags app富文章管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetAppMgmtRichArticleDetailReq true "请求参数"
// @Success 200 {object} pb.GetAppMgmtRichArticleDetailResp "响应数据"
// @Router /ms/get/richArticle/detail [post]
func (r *AppMgrHandler) getRichArticleDetail(ctx *gin.Context) {
	in := &pb.GetAppMgmtRichArticleDetailReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().GetAppMgmtRichArticleDetail(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// addRichArticle 新增app富文章
// @Summary 新增app富文章
// @Description 使用此接口新增app富文章
// @Tags app富文章管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.AddAppMgmtRichArticleReq true "请求参数"
// @Success 200 {object} pb.AddAppMgmtRichArticleResp "响应数据"
// @Router /ms/add/richArticle [post]
func (r *AppMgrHandler) addRichArticle(ctx *gin.Context) {
	in := &pb.AddAppMgmtRichArticleReq{}
	if err := ctx.ShouldBind(in); err != nil {
		logx.Errorf("addRichArticle ctx.ShouldBind(in) err:%v", err)
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().AddAppMgmtRichArticle(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// updateRichArticle 更新app富文章
// @Summary 更新app富文章
// @Description 使用此接口更新app富文章
// @Tags app富文章管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.UpdateAppMgmtRichArticleReq true "请求参数"
// @Success 200 {object} pb.UpdateAppMgmtRichArticleResp "响应数据"
// @Router /ms/update/richArticle [post]
func (r *AppMgrHandler) updateRichArticle(ctx *gin.Context) {
	in := &pb.UpdateAppMgmtRichArticleReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().UpdateAppMgmtRichArticle(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// deleteRichArticle 删除app富文章
// @Summary 删除app富文章
// @Description 使用此接口删除app富文章
// @Tags app富文章管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.DeleteAppMgmtRichArticleReq true "请求参数"
// @Success 200 {object} pb.DeleteAppMgmtRichArticleResp "响应数据"
// @Router /ms/delete/richArticle [post]
func (r *AppMgrHandler) deleteRichArticle(ctx *gin.Context) {
	in := &pb.DeleteAppMgmtRichArticleReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().DeleteAppMgmtRichArticle(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}
