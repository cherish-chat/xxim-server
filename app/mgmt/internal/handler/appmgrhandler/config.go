package appmgrhandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
)

// getAllConfigList 获取全部配置列表
// @Summary 获取全部配置列表
// @Description 使用此接口获取全部配置列表
// @Tags app管理配置管理相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetAllAppMgmtConfigReq true "请求参数"
// @Success 200 {object} pb.GetAllAppMgmtConfigResp "响应数据"
// @Router /appmgmt/get/config/list/all [post]
func (r *AppMgrHandler) getAllConfigList(ctx *gin.Context) {
	in := &pb.GetAllAppMgmtConfigReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().GetAllAppMgmtConfig(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// updateAllConfigList 更新全部配置列表
// @Summary 更新全部配置列表
// @Description 使用此接口更新全部配置列表
// @Tags app管理配置管理相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.UpdateAppMgmtConfigReq true "请求参数"
// @Success 200 {object} pb.UpdateAppMgmtConfigResp "响应数据"
// @Router /appmgmt/update/config [post]
func (r *AppMgrHandler) updateAllConfigList(ctx *gin.Context) {
	in := &pb.UpdateAppMgmtConfigReq{}
	if err := ctx.ShouldBind(in); err != nil {
		logx.Errorf("updateAllConfigList err: %v", err)
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().UpdateAppMgmtConfig(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}
