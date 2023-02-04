package appmgrhandler

import "github.com/gin-gonic/gin"

// getAllConfigList 获取全部配置列表
// @Summary 获取全部配置列表
// @Description 使用此接口获取全部配置列表
// @Tags app管理配置管理相关接口
// @Accept application/json
// @Produce application/json
// @Param Token header string true "用户令牌"
// @Param UserId header string true "用户ID"
// @Param object body pb.GetAllMSIpWhiteListReq true "请求参数"
// @Success 200 {object} pb.GetAllMSIpWhiteListResp "响应数据"
// @Router /ms/get/ipwhitelist/list/all [post]
func (r *AppMgrHandler) getAllConfigList(ctx *gin.Context) {

}

func (r *AppMgrHandler) updateAllConfigList(ctx *gin.Context) {

}
