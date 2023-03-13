package grouphandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
)

// insertMemberZombie 向群里插入僵尸号
// @Summary 向群里插入僵尸号
// @Description 使用此接口向群里插入僵尸号
// @Tags 群聊成员管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param GroupId header string true "群聊ID"
// @Param object body pb.RandInsertZombieMemberReq true "请求参数"
// @Success 200 {object} pb.RandInsertZombieMemberResp "响应数据"
// @Router /group/insert/member/zombie [post]
func (r *GroupHandler) insertMemberZombie(ctx *gin.Context) {
	in := &pb.RandInsertZombieMemberReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.GroupService().RandInsertZombieMember(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// clearMemberZombie 清除群里的僵尸号
// @Summary 清除群里的僵尸号
// @Description 使用此接口清除群里的僵尸号
// @Tags 群聊成员管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param GroupId header string true "群聊ID"
// @Param object body pb.ClearZombieMemberReq true "请求参数"
// @Success 200 {object} pb.ClearZombieMemberResp "响应数据"
// @Router /group/clear/member/zombie [post]
func (r *GroupHandler) clearMemberZombie(ctx *gin.Context) {
	in := &pb.ClearZombieMemberReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.GroupService().ClearZombieMember(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}
