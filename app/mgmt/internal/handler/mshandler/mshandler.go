package mshandler

import (
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/gin-gonic/gin"
)

type MSHandler struct {
	svcCtx *svc.ServiceContext
}

func NewMSHandler(svcCtx *svc.ServiceContext) *MSHandler {
	return &MSHandler{svcCtx: svcCtx}
}

func (r *MSHandler) Register(g *gin.RouterGroup) {
	group := g.Group("/ms")
	// 登录
	group.POST("/login", r.login)
	// 定期检查健康
	group.POST("/health", r.health)
	// 管理菜单
	{
		// 获取全部菜单列表
		group.POST("/get/menu/list/all", r.getAllMenuList)
		// 获取我权限内的菜单列表
		group.POST("/get/menu/list", r.getMenuList)
		// 获取详情
		group.POST("/get/menu/detail", r.getMenuDetail)
		// 增加菜单
		group.POST("/add/menu", r.addMenu)
		// 修改菜单
		group.POST("/update/menu", r.updateMenu)
		// 批量删除菜单
		group.POST("/delete/menu/batch", r.deleteMenuBatch)
	}
	// 管理后端api path
	{
		// 获取全部api path列表
		group.POST("/get/apipath/list/all", r.getAllApiPathList)
		// 获取我权限内的api path列表
		group.POST("/get/apipath/list", r.getApiPathList)
		// 获取详情
		group.POST("/get/apipath/detail", r.getApiPathDetail)
		// 增加api path
		group.POST("/add/apipath", r.addApiPath)
		// 修改api path
		group.POST("/update/apipath", r.updateApiPath)
		// 批量删除api path
		group.POST("/delete/apipath/batch", r.deleteApiPathBatch)
	}
	// 管理角色
	{
		// 获取全部角色列表
		group.POST("/get/role/list/all", r.getAllRoleList)
		// 获取详情
		group.POST("/get/role/detail", r.getRoleDetail)
		// 增加角色
		group.POST("/add/role", r.addRole)
		// 修改角色
		group.POST("/update/role", r.updateRole)
		// 批量删除角色
		group.POST("/delete/role/batch", r.deleteRoleBatch)
	}
	// 角色绑定菜单
	{
		// 绑定角色菜单
		group.POST("/bind/role/menu/batch", r.bindRoleMenuBatch)
		// 解绑角色菜单
		group.POST("/unbind/role/menu/batch", r.unbindRoleMenuBatch)
	}
	// 角色绑定api path
	{
		// 绑定角色api path
		group.POST("/bind/role/apipath/batch", r.bindRoleApiPathBatch)
		// 解绑角色api path
		group.POST("/unbind/role/apipath/batch", r.unbindRoleApiPathBatch)
	}
	// 管理管理员
	{
		// 获取全部管理员列表
		group.POST("/get/admin/list/all", r.getAllAdminList)
		// 获取详情
		group.POST("/get/admin/detail", r.getAdminDetail)
		// 增加管理员
		group.POST("/add/admin", r.addAdmin)
		// 修改管理员
		group.POST("/update/admin", r.updateAdmin)
		// 批量删除管理员
		group.POST("/delete/admin/batch", r.deleteAdminBatch)
	}
	// 管理员绑定角色
	{
		// 绑定管理员角色
		group.POST("/bind/admin/role/batch", r.bindAdminRoleBatch)
		// 解绑管理员角色
		group.POST("/unbind/admin/role/batch", r.unbindAdminRoleBatch)
	}
	// ip白名单
	{
		// 获取全部ip白名单列表
		group.POST("/get/ipwhitelist/list/all", r.getAllIpWhiteList)
		// 获取详情
		group.POST("/get/ipwhitelist/detail", r.getIpWhiteListDetail)
		// 增加ip白名单
		group.POST("/add/ipwhitelist", r.addIpWhiteList)
		// 修改ip白名单
		group.POST("/update/ipwhitelist", r.updateIpWhiteList)
		// 批量删除ip白名单
		group.POST("/delete/ipwhitelist/batch", r.deleteIpWhiteListBatch)
	}
}
