package mshandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/handler"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/logic"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
)

// getAllOperationLogList 获取全部操作日志列表
// @Summary 获取全部操作日志列表
// @Description 使用此接口获取全部操作日志列表
// @Tags 管理系统操作日志相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetAllMSOperationLogReq true "请求参数"
// @Success 200 {object} pb.GetAllMSOperationLogResp "响应数据"
// @Router /ms/get/operationlog/list/all [post]
func (r *MSHandler) getAllOperationLogList(ctx *gin.Context) {
	in := &pb.GetAllMSOperationLogReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewGetAllMSOperationLogLogic(ctx, r.svcCtx).GetAllMSOperationLog(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// getOperationLogDetail 获取操作日志详情
// @Summary 获取操作日志详情
// @Description 使用此接口获取操作日志详情
// @Tags 管理系统操作日志相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetMSOperationLogDetailReq true "请求参数"
// @Success 200 {object} pb.GetMSOperationLogDetailResp "响应数据"
// @Router /ms/get/operationlog/detail [post]
func (r *MSHandler) getOperationLogDetail(ctx *gin.Context) {
	in := &pb.GetMSOperationLogDetailReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewGetMSOperationLogDetailLogic(ctx, r.svcCtx).GetMSOperationLogDetail(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}

// deleteOperationLogBatch 删除操作日志
// @Summary 删除操作日志
// @Description 使用此接口删除操作日志
// @Tags 管理系统操作日志相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.DeleteMSOperationLogReq true "请求参数"
// @Success 200 {object} pb.DeleteMSOperationLogResp "响应数据"
// @Router /ms/delete/operationlog [post]
func (r *MSHandler) deleteOperationLogBatch(ctx *gin.Context) {
	in := &pb.DeleteMSOperationLogReq{}
	if err := ctx.ShouldBind(in); err != nil {
		ctx.AbortWithStatus(400)
		return
	}
	out, err := logic.NewDeleteMSOperationLogLogic(ctx, r.svcCtx).DeleteMSOperationLog(in)
	if err != nil {
		ctx.AbortWithStatus(500)
		return
	}
	handler.ReturnOk(ctx, out)
}
