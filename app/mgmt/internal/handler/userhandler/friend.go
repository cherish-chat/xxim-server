package userhandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
)

// getAllFriend 好友列表
// @Summary 获取好友列表
// @Description 使用此接口获取好友列表
// @Tags 用户好友管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetFriendListByUserIdReq true "请求参数"
// @Success 200 {object} pb.GetFriendListByUserIdResp "响应数据"
// @Router /ms/get/friend/list/all [post]
func (r *UserHandler) getAllFriend(ctx *gin.Context) {
	in := &pb.GetFriendListByUserIdReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.RelationService().GetFriendListByUserId(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

type deleteFriendReq struct {
	CommonReq *pb.CommonReq `json:"commonReq"`
	UserId    string        `json:"userId"`
	FriendId  string        `json:"friendId"`
}

// deleteFriend 删除好友
// @Summary 删除好友
// @Description 使用此接口删除好友
// @Tags 用户好友管理
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body deleteFriendReq true "请求参数"
// @Success 200 {object} pb.DeleteFriendResp "响应数据"
// @Router /ms/delete/friend [post]
func (r *UserHandler) deleteFriend(ctx *gin.Context) {
	in := &deleteFriendReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.RelationService().DeleteFriend(ctx, &pb.DeleteFriendReq{
		CommonReq: &pb.CommonReq{UserId: in.UserId},
		UserId:    in.FriendId,
		Block:     false,
	})
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}
