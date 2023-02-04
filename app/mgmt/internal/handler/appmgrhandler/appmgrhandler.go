package appmgrhandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/gin-gonic/gin"
)

type AppMgrHandler struct {
	svcCtx *svc.ServiceContext
}

func NewAppMgrHandler(svcCtx *svc.ServiceContext) *AppMgrHandler {
	return &AppMgrHandler{svcCtx: svcCtx}
}

func (r *AppMgrHandler) Register(g *gin.RouterGroup) {
	group := g.Group("/appmgr") // app管理
	// app基础配置
	{
		// 获取全部app基础配置列表
		group.POST("/get/config/list/all", r.getAllConfigList)
		// 更新全部app基础配置列表
		group.POST("/update/config/list/all", r.updateAllConfigList)
	}
	// app版本管理
	{
		// 获取全部app版本列表
		group.POST("/get/version/list/all", r.getAllVersionList)
		// 获取app版本详情
		group.POST("/get/version/detail", r.getVersionDetail)
		// 增加app版本
		group.POST("/add/version", r.addVersion)
		// 更新app版本
		group.POST("/update/version", r.updateVersion)
		// 删除app版本
		group.POST("/delete/version", r.deleteVersion)
	}
	// app屏蔽词管理
	{
		// 获取全部app屏蔽词列表
		group.POST("/get/shieldword/list/all", r.getAllShieldWordList)
		// 获取app屏蔽词详情
		group.POST("/get/shieldword/detail", r.getShieldWordDetail)
		// 增加app屏蔽词
		group.POST("/add/shieldword", r.addShieldWord)
		// 更新app屏蔽词
		group.POST("/update/shieldword", r.updateShieldWord)
		// 删除app屏蔽词
		group.POST("/delete/shieldword", r.deleteShieldWord)
		// 批量增加app屏蔽词 使用\n分割
		group.POST("/add/shieldword/batch", r.addShieldWordBatch)
	}
	// app连接管理 （长连接）
	{
		// 获取全部app连接列表
		group.POST("/get/connection/list/all", r.getAllConnectionList)
		// 获取app连接详情
		group.POST("/get/connection/detail", r.getConnectionDetail)
		// 踢出app连接
		group.POST("/kickout/connection", r.kickoutConnection)
	}
	// VPN列表
	{
		// 获取全部VPN列表
		group.POST("/get/vpn/list/all", r.getAllVPNList)
		// 获取VPN详情
		group.POST("/get/vpn/detail", r.getVPNDetail)
		// 增加VPN
		group.POST("/add/vpn", r.addVPN)
		// 更新VPN
		group.POST("/update/vpn", r.updateVPN)
		// 删除VPN
		group.POST("/delete/vpn", r.deleteVPN)
	}
	// app内飘屏通知管理
	{
		// 获取全部app内飘屏通知列表
		group.POST("/get/notice/list/all", r.getAllNoticeList)
		// 获取app内飘屏通知详情
		group.POST("/get/notice/detail", r.getNoticeDetail)
		// 增加app内飘屏通知
		group.POST("/add/notice", r.addNotice)
		// 更新app内飘屏通知
		group.POST("/update/notice", r.updateNotice)
		// 删除app内飘屏通知
		group.POST("/delete/notice", r.deleteNotice)
	}
	// app公告管理
	{
		// 获取全部app公告列表
		group.POST("/get/announcement/list/all", r.getAllAnnouncementList)
		// 获取app公告详情
		group.POST("/get/announcement/detail", r.getAnnouncementDetail)
		// 增加app公告
		group.POST("/add/announcement", r.addAnnouncement)
		// 更新app公告
		group.POST("/update/announcement", r.updateAnnouncement)
		// 删除app公告
		group.POST("/delete/announcement", r.deleteAnnouncement)
	}
}
