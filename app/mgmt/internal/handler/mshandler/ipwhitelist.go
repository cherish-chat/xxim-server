package mshandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/logic"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
)

// getAllIpWhiteList 获取全部ip白名单列表
// @Summary 获取全部ip白名单列表
// @Description 使用此接口获取全部ip白名单列表
// @Tags 管理系统ip白名单相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param object body pb.GetAllMSIpWhiteListReq true "请求参数"
// @Success 200 {object} pb.GetAllMSIpWhiteListResp "响应数据"
// @Router /ms/get/ipwhitelist/list/all [post]
func (r *MSHandler) getAllIpWhiteList(ctx *gin.Context) {
	in := &pb.GetAllMSIpWhiteListReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewGetAllMSIpWhiteListLogic(ctx, r.svcCtx).GetAllMSIpWhiteList(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	ctx.JSON(200, out)
}

// getIpWhiteListDetail 获取ip白名单详情
// @Summary 获取ip白名单详情
// @Description 使用此接口获取ip白名单详情
// @Tags 管理系统ip白名单相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param object body pb.GetMSIpWhiteListDetailReq true "请求参数"
// @Success 200 {object} pb.GetMSIpWhiteListDetailResp "响应数据"
// @Router /ms/get/ipwhitelist/detail [post]
func (r *MSHandler) getIpWhiteListDetail(ctx *gin.Context) {
	in := &pb.GetMSIpWhiteListDetailReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewGetMSIpWhiteListDetailLogic(ctx, r.svcCtx).GetMSIpWhiteListDetail(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	ctx.JSON(200, out)
}

// addIpWhiteList 添加ip白名单
// @Summary 添加ip白名单
// @Description 使用此接口添加ip白名单
// @Tags 管理系统ip白名单相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param object body pb.AddMSIpWhiteListReq true "请求参数"
// @Success 200 {object} pb.AddMSIpWhiteListResp "响应数据"
// @Router /ms/add/ipwhitelist [post]
func (r *MSHandler) addIpWhiteList(ctx *gin.Context) {
	in := &pb.AddMSIpWhiteListReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewAddMSIpWhiteListLogic(ctx, r.svcCtx).AddMSIpWhiteList(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	ctx.JSON(200, out)
}

// updateIpWhiteList 更新ip白名单
// @Summary 更新ip白名单
// @Description 使用此接口更新ip白名单
// @Tags 管理系统ip白名单相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param object body pb.UpdateMSIpWhiteListReq true "请求参数"
// @Success 200 {object} pb.UpdateMSIpWhiteListResp "响应数据"
// @Router /ms/update/ipwhitelist [post]
func (r *MSHandler) updateIpWhiteList(ctx *gin.Context) {
	in := &pb.UpdateMSIpWhiteListReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewUpdateMSIpWhiteListLogic(ctx, r.svcCtx).UpdateMSIpWhiteList(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	ctx.JSON(200, out)
}

// deleteIpWhiteListBatch 删除ip白名单
// @Summary 删除ip白名单
// @Description 使用此接口删除ip白名单
// @Tags 管理系统ip白名单相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param object body pb.DeleteMSIpWhiteListReq true "请求参数"
// @Success 200 {object} pb.DeleteMSIpWhiteListResp "响应数据"
// @Router /ms/delete/ipwhitelist [post]
func (r *MSHandler) deleteIpWhiteListBatch(ctx *gin.Context) {
	in := &pb.DeleteMSIpWhiteListReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewDeleteMSIpWhiteListLogic(ctx, r.svcCtx).DeleteMSIpWhiteList(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	ctx.JSON(200, out)
}
