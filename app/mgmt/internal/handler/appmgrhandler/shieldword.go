package appmgrhandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
)

// getAllShieldWordList 获取全部app屏蔽词列表
// @Summary 获取全部app屏蔽词列表
// @Description 使用此接口获取全部app屏蔽词列表
// @Tags app屏蔽词管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetAllAppMgmtShieldWordReq true "请求参数"
// @Success 200 {object} pb.GetAllAppMgmtShieldWordResp "响应数据"
// @Router /ms/get/shieldword/list/all [post]
func (r *AppMgrHandler) getAllShieldWordList(ctx *gin.Context) {
	in := &pb.GetAllAppMgmtShieldWordReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().GetAllAppMgmtShieldWord(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// getShieldWordDetail 获取app屏蔽词详情
// @Summary 获取app屏蔽词详情
// @Description 使用此接口获取app屏蔽词详情
// @Tags app屏蔽词管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetAppMgmtShieldWordDetailReq true "请求参数"
// @Success 200 {object} pb.GetAppMgmtShieldWordDetailResp "响应数据"
// @Router /ms/get/shieldword/detail [post]
func (r *AppMgrHandler) getShieldWordDetail(ctx *gin.Context) {
	in := &pb.GetAppMgmtShieldWordDetailReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().GetAppMgmtShieldWordDetail(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// addShieldWord 新增app屏蔽词
// @Summary 新增app屏蔽词
// @Description 使用此接口新增app屏蔽词
// @Tags app屏蔽词管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.AddAppMgmtShieldWordReq true "请求参数"
// @Success 200 {object} pb.AddAppMgmtShieldWordResp "响应数据"
// @Router /ms/add/shieldword [post]
func (r *AppMgrHandler) addShieldWord(ctx *gin.Context) {
	in := &pb.AddAppMgmtShieldWordReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().AddAppMgmtShieldWord(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// updateShieldWord 更新app屏蔽词
// @Summary 更新app屏蔽词
// @Description 使用此接口更新app屏蔽词
// @Tags app屏蔽词管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.UpdateAppMgmtShieldWordReq true "请求参数"
// @Success 200 {object} pb.UpdateAppMgmtShieldWordResp "响应数据"
// @Router /ms/update/shieldword [post]
func (r *AppMgrHandler) updateShieldWord(ctx *gin.Context) {
	in := &pb.UpdateAppMgmtShieldWordReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().UpdateAppMgmtShieldWord(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// deleteShieldWord 删除app屏蔽词
// @Summary 删除app屏蔽词
// @Description 使用此接口删除app屏蔽词
// @Tags app屏蔽词管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.DeleteAppMgmtShieldWordReq true "请求参数"
// @Success 200 {object} pb.DeleteAppMgmtShieldWordResp "响应数据"
// @Router /ms/delete/shieldword [post]
func (r *AppMgrHandler) deleteShieldWord(ctx *gin.Context) {
	in := &pb.DeleteAppMgmtShieldWordReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().DeleteAppMgmtShieldWord(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}
