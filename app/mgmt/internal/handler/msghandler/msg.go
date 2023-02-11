package msghandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
)

// getAllMsg 获取全部消息模型列表
// @Summary 获取全部消息模型列表
// @Description 使用此接口获取全部消息模型列表
// @Tags 消息模型管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "消息令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetAllMsgListReq true "请求参数"
// @Success 200 {object} pb.GetAllMsgListResp "响应数据"
// @Router /ms/get/msg/list/all [post]
func (r *MsgHandler) getAllMsg(ctx *gin.Context) {
	in := &pb.GetAllMsgListReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.MsgService().GetAllMsgList(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}
