package mshandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/logic"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
)

// login 登录managersystem管理系统
// @Summary 登录managersystem管理系统
// @Description 必须是管理员才能登录
// @Tags 管理系统相关接口
// @Accept application/json
// @Produce application/json
// @Param object body pb.LoginMSReq true "请求参数"
// @Success 200 {object} pb.LoginMSResp "响应数据"
// @Router /ms/login [post]
func (r *MSHandler) login(ctx *gin.Context) {
	in := &pb.LoginMSReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewLoginMSLogic(ctx, r.svcCtx).LoginMS(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	if out.GetCommonResp().GetCode() != pb.CommonResp_Success {
		handler.ReturnError(ctx, out.GetCommonResp())
		return
	}
	handler.ReturnOk(ctx, out)
}

// loginCaptcha 登录managersystem管理系统
// @Summary 登录managersystem管理系统
// @Description 必须是管理员才能登录
// @Tags 管理系统相关接口
// @Accept application/json
// @Produce application/json
// @Param object body pb.LoginMSCaptchaReq true "请求参数"
// @Success 200 {object} pb.LoginMSCaptchaResp "响应数据"
// @Router /ms/login/captcha [post]
func (r *MSHandler) loginCaptcha(ctx *gin.Context) {
	in := &pb.LoginMSCaptchaReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewLoginMSCaptchaLogic(ctx, r.svcCtx).LoginMSCaptcha(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// health 检查服务是否健康
// @Summary 检查服务是否健康
// @Description 使用此接口检查服务是否健康
// @Tags 管理系统相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Success 200 {object} pb.HealthMSResp "响应数据"
// @Router /ms/health [post]
func (r *MSHandler) health(ctx *gin.Context) {
	in := &pb.CommonReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewHealthMSLogic(ctx, r.svcCtx).HealthMS(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// config 获取配置信息
// @Summary 获取配置信息
// @Description 获取配置信息
// @Tags 管理系统相关接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} pb.ConfigMSResp "响应数据"
// @Router /ms/config [post]
func (r *MSHandler) config(ctx *gin.Context) {
	in := &pb.CommonReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewConfigMSLogic(ctx, r.svcCtx).ConfigMS(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// stats 统计
// @Summary 统计
// @Description 统计
// @Tags 管理系统相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.StatsMSReq true "请求参数"
// @Success 200 {object} pb.StatsMSResp "响应数据"
// @Router /ms/stats [post]
func (r *MSHandler) stats(ctx *gin.Context) {
	in := &pb.StatsMSReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewStatsMSLogic(ctx, r.svcCtx).StatsMS(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// clearAllUser 清空所有用户
// @Summary 清空所有用户
// @Description 清空所有用户
// @Tags 管理系统相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.ClearAllUserReq true "请求参数"
// @Success 200 {object} pb.ClearAllUserResp "响应数据"
// @Router /ms/clear/all/user [post]
func (r *MSHandler) clearAllUser(ctx *gin.Context) {
	in := &pb.ClearAllUserReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewClearAllUserLogic(ctx, r.svcCtx).ClearAllUser(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}
