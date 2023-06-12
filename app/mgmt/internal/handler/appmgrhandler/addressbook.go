package appmgrhandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
)

// getAddressBook 获取通讯录
// @Summary 获取通讯录
// @Description 使用此接口获取通讯录
// @Tags app管理配置管理相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetAppAddressBookReq true "请求参数"
// @Success 200 {object} pb.GetAppAddressBookResp "响应数据"
// @Router /appmgmt/get/addressbook [post]
func (r *AppMgrHandler) getAddressBook(ctx *gin.Context) {
	in := &pb.GetAppAddressBookReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().GetAppAddressBook(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// updateAddressBook 更新通讯录
// @Summary 更新通讯录
// @Description 使用此接口更新通讯录
// @Tags app管理配置管理相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.UpdateAppAddressBookReq true "请求参数"
// @Success 200 {object} pb.UpdateAppAddressBookResp "响应数据"
// @Router /appmgmt/update/addressbook [post]
func (r *AppMgrHandler) updateAddressBook(ctx *gin.Context) {
	in := &pb.UpdateAppAddressBookReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := r.svcCtx.AppMgmtService().UpdateAppAddressBook(ctx, in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}
