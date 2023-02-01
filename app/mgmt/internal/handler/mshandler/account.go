package mshandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/logic"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
)

// login 登录managersystem管理系统
// @Summary 登录managersystem管理系统
// @Description 必须是管理员才能登录
// @Tags 账号相关接口
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
	ctx.JSON(200, out)
}

// health 检查服务是否健康
// @Summary 检查服务是否健康
// @Description 使用此接口检查服务是否健康
// @Tags 账号相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
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
	ctx.JSON(200, out)
}
