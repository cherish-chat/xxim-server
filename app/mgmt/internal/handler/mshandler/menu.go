package mshandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/logic"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
)

// getAllMenuList 获取所有菜单列表
// @Summary 获取所有菜单列表
// @Description 使用此接口获取所有菜单列表
// @Tags 管理系统菜单相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param object body pb.GetAllMSMenuListReq true "请求参数"
// @Success 200 {object} pb.GetAllMSMenuListResp "响应数据"
// @Router /ms/get/menu/list/all [post]
func (r *MSHandler) getAllMenuList(ctx *gin.Context) {
	in := &pb.GetAllMSMenuListReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewGetAllMSMenuListLogic(ctx, r.svcCtx).GetAllMSMenuList(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// getMenuList 获取我的菜单列表
// @Summary 获取我的菜单列表
// @Description 使用此接口获取我的菜单列表
// @Tags 管理系统菜单相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param object body pb.GetMyMSMenuListReq true "请求参数"
// @Success 200 {object} pb.GetMyMSMenuListResp "响应数据"
// @Router /ms/get/menu/list [post]
func (r *MSHandler) getMenuList(ctx *gin.Context) {
	in := &pb.GetMyMSMenuListReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewGetMyMSMenuListLogic(ctx, r.svcCtx).GetMyMSMenuList(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// getMenuDetail 获取菜单详情
// @Summary 获取菜单详情
// @Description 使用此接口获取菜单详情
// @Tags 管理系统菜单相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param object body pb.GetMSMenuDetailReq true "请求参数"
// @Success 200 {object} pb.GetMSMenuDetailResp "响应数据"
// @Router /ms/get/menu/detail [post]
func (r *MSHandler) getMenuDetail(ctx *gin.Context) {
	in := &pb.GetMSMenuDetailReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewGetMSMenuDetailLogic(ctx, r.svcCtx).GetMSMenuDetail(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// addMenu 新增菜单
// @Summary 新增菜单
// @Description 使用此接口新增菜单
// @Tags 管理系统菜单相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param object body pb.AddMSMenuReq true "请求参数"
// @Success 200 {object} pb.AddMSMenuResp "响应数据"
// @Router /ms/add/menu [post]
func (r *MSHandler) addMenu(ctx *gin.Context) {
	in := &pb.AddMSMenuReq{}
	if err := ctx.ShouldBind(in); err != nil {
		logx.Errorf("addMenu err: %v", err)
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewAddMSMenuLogic(ctx, r.svcCtx).AddMSMenu(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// updateMenu 更新菜单
// @Summary 更新菜单
// @Description 使用此接口更新菜单
// @Tags 管理系统菜单相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param object body pb.UpdateMSMenuReq true "请求参数"
// @Success 200 {object} pb.UpdateMSMenuResp "响应数据"
// @Router /ms/update/menu [post]
func (r *MSHandler) updateMenu(ctx *gin.Context) {
	in := &pb.UpdateMSMenuReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewUpdateMSMenuLogic(ctx, r.svcCtx).UpdateMSMenu(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// deleteMenuBatch 删除菜单
// @Summary 删除菜单
// @Description 使用此接口删除菜单
// @Tags 管理系统菜单相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param object body pb.DeleteMSMenuReq true "请求参数"
// @Success 200 {object} pb.DeleteMSMenuResp "响应数据"
// @Router /ms/delete/menu [post]
func (r *MSHandler) deleteMenuBatch(ctx *gin.Context) {
	in := &pb.DeleteMSMenuReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewDeleteMSMenuLogic(ctx, r.svcCtx).DeleteMSMenu(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}
