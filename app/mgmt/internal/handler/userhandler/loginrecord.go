package userhandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
)

// getAllLoginRecord 获取全部用户登录记录列表
// @Summary 获取全部用户登录记录列表
// @Description 使用此接口获取全部用户登录记录列表
// @Tags 用户登录记录管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetAllLoginRecordReq true "请求参数"
// @Success 200 {object} pb.GetAllLoginRecordResp "响应数据"
// @Router /ms/get/loginrecord/list/all [post]
func (r *UserHandler) getAllLoginRecord(ctx *gin.Context) {
	in := &pb.GetAllLoginRecordReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.UserService().GetAllLoginRecord(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}
