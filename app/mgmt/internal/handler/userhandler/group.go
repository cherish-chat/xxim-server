package userhandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/gin-gonic/gin"
)

// getAllGroup 群列表
// @Summary 获取群列表
// @Description 使用此接口获取群列表
// @Tags 用户群管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetGroupListByUserIdReq true "请求参数"
// @Success 200 {object} pb.GetGroupListByUserIdResp "响应数据"
// @Router /ms/get/group/list/all [post]
func (r *UserHandler) getAllGroup(ctx *gin.Context) {
	in := &pb.GetGroupListByUserIdReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.GroupService().GetGroupListByUserId(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

type deleteGroupReq struct {
	CommonReq *pb.CommonReq `json:"commonReq"`
	UserId    string        `json:"userId"`
	GroupId   string        `json:"groupId"`
}

// deleteGroup 删除群
// @Summary 删除群
// @Description 使用此接口删除群
// @Tags 用户群管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body deleteGroupReq true "请求参数"
// @Success 200 {object} pb.KickGroupMemberResp "响应数据"
func (r *UserHandler) deleteGroup(ctx *gin.Context) {
	in := &deleteGroupReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	// 获取群主
	memberList, err := r.svcCtx.GroupService().GetGroupMemberList(ctx, &pb.GetGroupMemberListReq{
		CommonReq: &pb.CommonReq{UserId: in.UserId},
		GroupId:   in.GroupId,
		Page: &pb.Page{
			Page: 1,
			Size: 1,
		},
		Filter: &pb.GetGroupMemberListReq_GetGroupMemberListFilter{
			OnlyOwner: utils.AnyPtr(true),
		},
		Opt: &pb.GetGroupMemberListReq_GetGroupMemberListOpt{OnlyId: utils.AnyPtr(true)},
	})
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	if len(memberList.GroupMemberList) == 0 {
		ctx.AbortWithStatus(500)
		return
	}
	out, err := r.svcCtx.GroupService().KickGroupMember(ctx, &pb.KickGroupMemberReq{
		CommonReq: &pb.CommonReq{UserId: memberList.GroupMemberList[0].MemberId},
		GroupId:   in.GroupId,
		MemberId:  in.UserId,
	})
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}
