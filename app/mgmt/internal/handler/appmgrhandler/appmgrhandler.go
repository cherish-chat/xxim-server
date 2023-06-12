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
	group := g.Group("/appmgmt") // app管理
	// app基础配置
	{
		// 获取全部app基础配置列表
		group.POST("/get/config/list/all", r.getAllConfigList)
		// 更新全部app基础配置列表
		group.POST("/update/config", r.updateAllConfigList)
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
	}
	// app连接管理 （长连接）
	{
		// 获取全部app连接列表
		group.POST("/get/connection/list/all", r.getAllConnectionList)
		// 踢出app连接
		group.POST("/kickout/connection", r.kickoutConnection)
	}
	// VPN列表
	{
		// 获取全部VPN列表
		group.POST("/get/vpn/list/all", r.getAllVpnList)
		// 获取VPN详情
		group.POST("/get/vpn/detail", r.getVpnDetail)
		// 增加VPN
		group.POST("/add/vpn", r.addVpn)
		// 更新VPN
		group.POST("/update/vpn", r.updateVpn)
		// 删除VPN
		group.POST("/delete/vpn", r.deleteVpn)
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
	// app内富文章
	{
		// 获取全部app内富文章列表
		group.POST("/get/richArticle/list/all", r.getAllRichArticleList)
		// 获取app内富文章详情
		group.POST("/get/richArticle/detail", r.getRichArticleDetail)
		// 增加app内富文章
		group.POST("/add/richArticle", r.addRichArticle)
		// 更新app内富文章
		group.POST("/update/richArticle", r.updateRichArticle)
		// 删除app内富文章
		group.POST("/delete/richArticle", r.deleteRichArticle)
	}
	// app内表情组管理
	{
		// 获取全部app内表情组列表
		group.POST("/get/emojigroup/list/all", r.getAllEmojiGroupList)
		// 获取app内表情组详情
		group.POST("/get/emojigroup/detail", r.getEmojiGroupDetail)
		// 更新app内表情组
		group.POST("/update/emojigroup", r.updateEmojiGroup)
		// 删除app内表情组和表情
		group.POST("/delete/emojigroup", r.deleteEmojiGroup)
	}
	// app内表情管理
	{
		// 获取全部app内表情列表
		group.POST("/get/emoji/list/all", r.getAllEmojiList)
		// 获取app内表情详情
		group.POST("/get/emoji/detail", r.getEmojiDetail)
		// 增加app内表情
		group.POST("/add/emoji", r.addEmoji)
		// 更新app内表情
		group.POST("/update/emoji", r.updateEmoji)
		// 删除app内表情
		group.POST("/delete/emoji", r.deleteEmoji)
	}
	// 外链列表
	{
		// 获取全部VPN列表
		group.POST("/get/link/list/all", r.getAllLinkList)
		// 获取VPN详情
		group.POST("/get/link/detail", r.getLinkDetail)
		// 增加VPN
		group.POST("/add/link", r.addLink)
		// 更新VPN
		group.POST("/update/link", r.updateLink)
		// 删除VPN
		group.POST("/delete/link", r.deleteLink)
	}
	// app通讯录
	{
		// 获取app通讯录
		group.POST("/get/addressbook", r.getAddressBook)
		// 更新app通讯录
		group.POST("/update/addressbook", r.updateAddressBook)
	}
}
