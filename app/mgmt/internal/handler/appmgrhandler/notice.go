package appmgrhandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
)

// getAllNoticeList 获取全部app公告列表
// @Summary 获取全部app公告列表
// @Description 使用此接口获取全部app公告列表
// @Tags app公告管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetAllAppMgmtNoticeReq true "请求参数"
// @Success 200 {object} pb.GetAllAppMgmtNoticeResp "响应数据"
// @Router /ms/get/notice/list/all [post]
func (r *AppMgrHandler) getAllNoticeList(ctx *gin.Context) {
	in := &pb.GetAllAppMgmtNoticeReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().GetAllAppMgmtNotice(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// getNoticeDetail 获取app公告详情
// @Summary 获取app公告详情
// @Description 使用此接口获取app公告详情
// @Tags app公告管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetAppMgmtNoticeDetailReq true "请求参数"
// @Success 200 {object} pb.GetAppMgmtNoticeDetailResp "响应数据"
// @Router /ms/get/notice/detail [post]
func (r *AppMgrHandler) getNoticeDetail(ctx *gin.Context) {
	in := &pb.GetAppMgmtNoticeDetailReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().GetAppMgmtNoticeDetail(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// addNotice 新增app公告
// @Summary 新增app公告
// @Description 使用此接口新增app公告
// @Tags app公告管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.AddAppMgmtNoticeReq true "请求参数"
// @Success 200 {object} pb.AddAppMgmtNoticeResp "响应数据"
// @Router /ms/add/notice [post]
func (r *AppMgrHandler) addNotice(ctx *gin.Context) {
	in := &pb.AddAppMgmtNoticeReq{}
	if err := ctx.ShouldBind(in); err != nil {
		logx.Errorf("addNotice ctx.ShouldBind(in) err:%v", err)
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().AddAppMgmtNotice(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// updateNotice 更新app公告
// @Summary 更新app公告
// @Description 使用此接口更新app公告
// @Tags app公告管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.UpdateAppMgmtNoticeReq true "请求参数"
// @Success 200 {object} pb.UpdateAppMgmtNoticeResp "响应数据"
// @Router /ms/update/notice [post]
func (r *AppMgrHandler) updateNotice(ctx *gin.Context) {
	in := &pb.UpdateAppMgmtNoticeReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().UpdateAppMgmtNotice(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// deleteNotice 删除app公告
// @Summary 删除app公告
// @Description 使用此接口删除app公告
// @Tags app公告管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.DeleteAppMgmtNoticeReq true "请求参数"
// @Success 200 {object} pb.DeleteAppMgmtNoticeResp "响应数据"
// @Router /ms/delete/notice [post]
func (r *AppMgrHandler) deleteNotice(ctx *gin.Context) {
	in := &pb.DeleteAppMgmtNoticeReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().DeleteAppMgmtNotice(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}
