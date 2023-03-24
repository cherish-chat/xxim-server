package userhandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svcCtx *svc.ServiceContext
}

func NewUserHandler(svcCtx *svc.ServiceContext) *UserHandler {
	return &UserHandler{svcCtx: svcCtx}
}

func (r *UserHandler) Register(g *gin.RouterGroup) {
	group := g.Group("/usermgmt") // app管理
	{
		// DefaultConv 用户默认会话
		// 列表
		group.POST("/get/defaultconv/list/all", r.getAllDefaultConv)
		// 详情
		group.POST("/get/defaultconv/detail", r.getDefaultConvDetail)
		// 增加
		group.POST("/add/defaultconv", r.addDefaultConv)
		// 更新
		group.POST("/update/defaultconv", r.updateDefaultConv)
		// 删除
		group.POST("/delete/defaultconv", r.deleteDefaultConv)
	}
	{
		// InvitationCode 邀请码
		// 列表
		group.POST("/get/invitationcode/list/all", r.getAllInvitationCode)
		// 详情
		group.POST("/get/invitationcode/detail", r.getInvitationCodeDetail)
		// 增加
		group.POST("/add/invitationcode", r.addInvitationCode)
		// 更新
		group.POST("/update/invitationcode", r.updateInvitationCode)
		// 删除
		group.POST("/delete/invitationcode", r.deleteInvitationCode)
	}
	{
		// IpBlackList IP黑名单
		// 列表
		group.POST("/get/ipblacklist/list/all", r.getAllIpBlackList)
		// 详情
		group.POST("/get/ipblacklist/detail", r.getIpBlackListDetail)
		// 增加
		group.POST("/add/ipblacklist", r.addIpBlackList)
		// 更新
		group.POST("/update/ipblacklist", r.updateIpBlackList)
		// 删除
		group.POST("/delete/ipblacklist", r.deleteIpBlackList)
	}
	{
		// IpWhiteList IP白名单
		// 列表
		group.POST("/get/ipwhitelist/list/all", r.getAllIpWhiteList)
		// 详情
		group.POST("/get/ipwhitelist/detail", r.getIpWhiteListDetail)
		// 增加
		group.POST("/add/ipwhitelist", r.addIpWhiteList)
		// 更新
		group.POST("/update/ipwhitelist", r.updateIpWhiteList)
		// 删除
		group.POST("/delete/ipwhitelist", r.deleteIpWhiteList)
	}
	{
		// Model 用户
		// 列表
		group.POST("/get/model/list/all", r.getAllModel)
		// 详情
		group.POST("/get/model/detail", r.getModelDetail)
		// 增加
		group.POST("/add/model", r.addModel)
		// 更新
		group.POST("/update/model", r.updateModel)
		// 删除
		group.POST("/delete/model", r.deleteModel)
		// 切换状态
		group.POST("/switch/model", r.switchModel)
		// 批量创建僵尸号
		group.POST("/batch/create/zombie", r.batchCreateZombie)
	}
	{
		// LoginRecord 登录记录
		// 列表
		group.POST("/get/loginrecord/list/all", r.getAllLoginRecord)
	}
	{
		// Friend 好友
		// 列表
		group.POST("/get/friend/list/all", r.getAllFriend)
		// 删除
		group.POST("/delete/friend", r.deleteFriend)
	}
	{
		// Group 群组
		// 列表
		group.POST("/get/group/list/all", r.getAllGroup)
		// 删除
		group.POST("/delete/group", r.deleteGroup)
	}
}
