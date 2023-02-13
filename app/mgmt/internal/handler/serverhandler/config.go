package serverhandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/logic"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
)

// configAll 获取服务端的所有配置信息
// @Summary 获取服务端的所有配置信息
// @Description 使用此接口获取服务端的所有配置信息, 比如redis的配置信息
// @Tags 服务端相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Success 200 {object} pb.GetServerAllConfigResp "所有配置信息"
// @Router /server/get/config/all [post]
func (r *ServerHandler) configAll(ctx *gin.Context) {
	in := &pb.GetServerAllConfigReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	config, err := logic.NewGetServerAllConfigLogic(ctx, r.svcCtx).GetServerAllConfig(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, config)
}

// updateConfig 更新服务端的配置信息
// @Summary 更新服务端的配置信息
// @Description 使用此接口更新服务端的配置信息, 比如redis的配置信息
// @Tags 服务端相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param config body pb.UpdateServerConfigReq true "配置信息"
// @Success 200 {object} pb.UpdateServerConfigResp "更新结果"
// @Router /server/update/config [post]
func (r *ServerHandler) updateConfig(ctx *gin.Context) {
	in := &pb.UpdateServerConfigReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	resp, err := logic.NewUpdateServerConfigLogic(ctx, r.svcCtx).UpdateServerConfig(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, resp)
}

// appLine 获取app线路配置
// @Summary 获取app线路配置
// @Description 使用此接口获取app线路配置
// @Tags 服务端相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Success 200 {object} pb.GetAppLineConfigResp "app线路配置"
// @Router /server/get/app/line [post]
func (r *ServerHandler) appLine(ctx *gin.Context) {
	in := &pb.GetAppLineConfigReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	resp, err := logic.NewGetAppLineConfigLogic(ctx, r.svcCtx).GetAppLineConfig(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, resp)
}

// updateAppLine 更新app线路配置
// @Summary 更新app线路配置
// @Description 使用此接口更新app线路配置
// @Tags 服务端相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param config body pb.UpdateAppLineConfigReq true "app线路配置"
// @Success 200 {object} pb.UpdateAppLineConfigResp "更新结果"
// @Router /server/update/app/line [post]
func (r *ServerHandler) updateAppLine(ctx *gin.Context) {
	in := &pb.UpdateAppLineConfigReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	resp, err := logic.NewUpdateAppLineConfigLogic(ctx, r.svcCtx).UpdateAppLineConfig(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, resp)
}
