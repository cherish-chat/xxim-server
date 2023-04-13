package mshandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/logic"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
)

// getAllLuaConfig 获取全部lua脚本列表
// @Summary 获取全部lua脚本列表
// @Description 使用此接口获取全部lua脚本列表
// @Tags 管理系统lua脚本相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetAllMSLuaConfigReq true "请求参数"
// @Success 200 {object} pb.GetAllMSLuaConfigResp "响应数据"
// @Router /ms/get/luaconfig/list/all [post]
func (r *MSHandler) getAllLuaConfig(ctx *gin.Context) {
	in := &pb.GetAllMSLuaConfigReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewGetAllMSLuaConfigLogic(ctx, r.svcCtx).GetAllMSLuaConfig(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// getLuaConfigDetail 获取lua脚本详情
// @Summary 获取lua脚本详情
// @Description 使用此接口获取lua脚本详情
// @Tags 管理系统lua脚本相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetMSLuaConfigReq true "请求参数"
// @Success 200 {object} pb.GetMSLuaConfigResp "响应数据"
// @Router /ms/get/luaconfig/detail [post]
func (r *MSHandler) getLuaConfigDetail(ctx *gin.Context) {
	in := &pb.GetMSLuaConfigReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewGetMSLuaConfigDetailLogic(ctx, r.svcCtx).GetMSLuaConfigDetail(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// addLuaConfig 添加lua脚本
// @Summary 添加lua脚本
// @Description 使用此接口添加lua脚本
// @Tags 管理系统lua脚本相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.AddMSLuaConfigReq true "请求参数"
// @Success 200 {object} pb.AddMSLuaConfigResp "响应数据"
// @Router /ms/add/luaconfig [post]
func (r *MSHandler) addLuaConfig(ctx *gin.Context) {
	in := &pb.AddMSLuaConfigReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewAddMSLuaConfigLogic(ctx, r.svcCtx).AddMSLuaConfig(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// updateLuaConfig 更新lua脚本
// @Summary 更新lua脚本
// @Description 使用此接口更新lua脚本
// @Tags 管理系统lua脚本相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.UpdateMSLuaConfigReq true "请求参数"
// @Success 200 {object} pb.UpdateMSLuaConfigResp "响应数据"
// @Router /ms/update/luaconfig [post]
func (r *MSHandler) updateLuaConfig(ctx *gin.Context) {
	in := &pb.UpdateMSLuaConfigReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewUpdateMSLuaConfigLogic(ctx, r.svcCtx).UpdateMSLuaConfig(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// deleteLuaConfigBatch 删除lua脚本
// @Summary 删除lua脚本
// @Description 使用此接口删除lua脚本
// @Tags 管理系统lua脚本相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.DeleteMSLuaConfigReq true "请求参数"
// @Success 200 {object} pb.DeleteMSLuaConfigResp "响应数据"
// @Router /ms/delete/luaconfig [post]
func (r *MSHandler) deleteLuaConfigBatch(ctx *gin.Context) {
	in := &pb.DeleteMSLuaConfigReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewDeleteMSLuaConfigLogic(ctx, r.svcCtx).DeleteMSLuaConfig(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}
