package mgmtmodel

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"gorm.io/gorm"
	"time"
)

type Menu struct {
	Id         string `gorm:"column:id;primarykey;comment:'主键';"`
	Pid        string `gorm:"column:pid;not null;default:0;comment:'上级菜单'"`
	MenuType   string `gorm:"column:menuType;not null;default:'';comment:'权限类型: M=目录，C=菜单，A=按钮''"`
	MenuName   string `gorm:"column:menuName;not null;default:'';comment:'菜单名称'"`
	MenuIcon   string `gorm:"column:menuIcon;not null;default:'';comment:'菜单图标'"`
	MenuSort   int32  `gorm:"column:menuSort;not null;default:0;comment:'菜单排序'"`
	Perms      string `gorm:"column:perms;not null;default:'';comment:'权限标识'"`
	Paths      string `gorm:"column:paths;not null;default:'';comment:'路由地址'"`
	Component  string `gorm:"column:component;not null;default:'';comment:'前端组件'"`
	Selected   string `gorm:"column:selected;not null;default:'';comment:'选中路径'"`
	Params     string `gorm:"column:params;not null;default:'';comment:'路由参数'"`
	IsCache    bool   `gorm:"column:isCache;not null;default:0;comment:'是否缓存: 0=否, 1=是''"`
	IsShow     bool   `gorm:"column:isShow;not null;default:1;comment:'是否显示: 0=否, 1=是'"`
	IsDisable  bool   `gorm:"column:isDisable;not null;default:0;comment:'是否禁用: 0=否, 1=是'"`
	CreateTime int64  `gorm:"column:createTime;not null;comment:'创建时间'"`
	UpdateTime int64  `gorm:"column:updateTime;not null;comment:'更新时间'"`
}

func (m *Menu) TableName() string {
	return MGMT_TABLE_PREFIX + "menu"
}

func (m *Menu) ToPb() *pb.MSMenu {
	return &pb.MSMenu{
		Id:           m.Id,
		CreatedAt:    m.CreateTime,
		CreatedAtStr: utils.TimeFormat(m.CreateTime),
		Pid:          m.Pid,
		MenuType:     m.MenuType,
		MenuName:     m.MenuName,
		MenuIcon:     m.MenuIcon,
		MenuSort:     m.MenuSort,
		Perms:        m.Perms,
		Paths:        m.Paths,
		Component:    m.Component,
		Selected:     m.Selected,
		Params:       m.Params,
		IsCache:      m.IsCache,
		IsShow:       m.IsShow,
		IsDisable:    m.IsDisable,
		Children:     make([]*pb.MSMenu, 0),
		UpdatedAt:    m.UpdateTime,
		UpdatedAtStr: utils.TimeFormat(m.UpdateTime),
	}
}

func genMenu(in []*pb.MSMenu) []*Menu {
	var menus []*Menu
	for _, v := range in {
		menus = append(menus, &Menu{
			Id:        v.Id,
			Pid:       v.Pid,
			MenuType:  v.MenuType,
			MenuName:  v.MenuName,
			MenuIcon:  v.MenuIcon,
			MenuSort:  v.MenuSort,
			Perms:     v.Perms,
			Paths:     v.Paths,
			Component: v.Component,
			Selected:  v.Selected,
			Params:    v.Params,
			IsCache:   v.IsCache,
			IsShow:    v.IsShow,
			IsDisable: v.IsDisable,
		})
		if len(v.Children) > 0 {
			menus = append(menus, genMenu(v.Children)...)
		}
	}
	return menus
}
func genMenu1(
	id string,
	mame string,
	icon string,
	sort int32,
	path string,
	children ...*pb.MSMenu,
) *pb.MSMenu {
	return &pb.MSMenu{
		Id:           id,
		UpdatedAtStr: time.Now().Format("2006-01-02 15:04:05"),
		UpdatedBy:    "system",
		Pid:          "",
		MenuType:     "M",
		MenuName:     mame,
		MenuIcon:     icon,
		MenuSort:     sort,
		Perms:        "",
		Paths:        path,
		Component:    "",
		Selected:     "",
		Params:       "",
		IsCache:      false,
		IsShow:       true,
		IsDisable:    false,
		Children:     children,
	}
}
func genMenu2(
	id string,
	pid string,
	mame string,
	icon string,
	sort int32,
	perms string,
	path string,
	component string,
	children ...*pb.MSMenu,
) *pb.MSMenu {
	return &pb.MSMenu{
		Id:           id,
		UpdatedAtStr: time.Now().Format("2006-01-02 15:04:05"),
		UpdatedBy:    "system",
		Pid:          pid,
		MenuType:     "C",
		MenuName:     mame,
		MenuIcon:     icon,
		MenuSort:     sort,
		Perms:        perms,
		Paths:        path,
		Component:    component,
		Selected:     "",
		Params:       "",
		IsCache:      true,
		IsShow:       true,
		IsDisable:    false,
		Children:     children,
	}
}
func genMenu3(
	id string,
	pid string,
	mame string,
	perms string,
) *pb.MSMenu {
	return &pb.MSMenu{
		Id:           id,
		UpdatedAtStr: time.Now().Format("2006-01-02 15:04:05"),
		UpdatedBy:    "system",
		Pid:          pid,
		MenuType:     "A",
		MenuName:     mame,
		MenuIcon:     "",
		MenuSort:     0,
		Perms:        perms,
		Paths:        "",
		Component:    "",
		Selected:     "",
		Params:       "",
		IsCache:      false,
		IsShow:       true,
		IsDisable:    false,
		Children:     nil,
	}
}
func initMenu(tx *gorm.DB) {
	menus := genMenu([]*pb.MSMenu{
		// INSERT INTO `la_system_auth_menu` VALUES (1, 0, 'C', '工作台', 'el-icon-Monitor', 50, 'index:console', 'workbench', 'workbench/index', '', '', 1, 1, 0, 1650341765, 1668672757);
		{
			Id:           "1",
			CreatedAt:    time.Now().UnixMilli(),
			CreatedAtStr: time.Now().Format("2006-01-02 15:04:05"),
			CreatedBy:    "system",
			UpdatedAt:    time.Now().UnixMilli(),
			UpdatedAtStr: time.Now().Format("2006-01-02 15:04:05"),
			UpdatedBy:    "system",
			Pid:          "",
			MenuType:     "C",
			MenuName:     "工作台",
			MenuIcon:     "el-icon-Monitor",
			MenuSort:     50,
			Perms:        "index:console",
			Paths:        "workbench",
			Component:    "workbench/index",
			Selected:     "",
			Params:       "",
			IsCache:      true,
			IsShow:       true,
			IsDisable:    false,
			Children:     make([]*pb.MSMenu, 0),
		},
		genMenu1("100", "权限管理", "el-icon-Lock", 44, "permission",
			genMenu2("101", "100", "管理员", "local-icon-wode", 0,
				"system:admin:list", "admin", "permission/admin/index",
				genMenu3("102", "101", "管理员详情", "system:admin:detail"),
				genMenu3("103", "101", "管理员新增", "system:admin:add"),
				genMenu3("104", "101", "管理员编辑", "system:admin:edit"),
				genMenu3("105", "101", "管理员删除", "system:admin:del"),
				genMenu3("106", "101", "管理员状态", "system:admin:disable"),
			),
			genMenu2("110", "100", "角色管理", "el-icon-Female", 0,
				"system:role:list", "role", "permission/role/index",
				genMenu3("111", "110", "角色详情", "system:role:detail"),
				genMenu3("112", "110", "角色新增", "system:role:add"),
				genMenu3("113", "110", "角色编辑", "system:role:edit"),
				genMenu3("114", "110", "角色删除", "system:role:del"),
			),
			genMenu2("120", "100", "菜单管理", "el-icon-Operation", 0,
				"system:menu:list", "menu", "permission/menu/index",
				genMenu3("121", "120", "菜单详情", "system:menu:detail"),
				genMenu3("122", "120", "菜单新增", "system:menu:add"),
				genMenu3("123", "120", "菜单编辑", "system:menu:edit"),
				genMenu3("124", "120", "菜单删除", "system:menu:del"),
			),
			genMenu2("130", "100", "api管理", "el-icon-Key", 0,
				"system:apipath:list", "apipath", "permission/apipath/index",
				genMenu3("131", "130", "api详情", "system:apipath:detail"),
				genMenu3("132", "130", "api新增", "system:apipath:add"),
				genMenu3("133", "130", "api编辑", "system:apipath:edit"),
				genMenu3("134", "130", "api删除", "system:apipath:del"),
			),
			genMenu2("140", "100", "ip白名单", "el-icon-List", 0,
				"system:ipwhitelist:list", "ipwhitelist", "permission/ipwhitelist/index",
				genMenu3("141", "140", "白名单详情", "system:ipwhitelist:detail"),
				genMenu3("142", "140", "白名单新增", "system:ipwhitelist:add"),
				genMenu3("143", "140", "白名单编辑", "system:ipwhitelist:edit"),
				genMenu3("144", "140", "白名单删除", "system:ipwhitelist:del"),
			),
			genMenu2("150", "100", "操作日志", "el-icon-Notebook", 0,
				"system:operationlog:list", "operationlog", "permission/operationlog/index",
				genMenu3("151", "150", "操作日志详情", "system:operationlog:detail"),
			),
			genMenu2("160", "100", "登录日志", "el-icon-CameraFilled", 0,
				"system:loginlog:list", "loginlog", "permission/loginlog/index",
			),
		),
		genMenu1("200", "运维管理", "local-icon-KMSguanli", 49, "devops",
			// config 配置管理
			genMenu2("201", "200", "配置管理", "el-icon-Setting", 0,
				"", "devops/config", "devops/config/index",
				genMenu3("202", "201", "更新配置", "devops:config:edit"),
			),
			// appline app线路
			genMenu2("211", "200", "app线路", "el-icon-Connection", 0,
				"", "devops/appline", "devops/appline/index",
				genMenu3("212", "211", "编辑线路", "devops:appline:edit"),
			),
		),
		genMenu1("300", "app管理", "el-icon-Apple", 48, "app",
			genMenu2("301", "300", "配置管理", "el-icon-Setting", 0,
				"", "app/config", "app/config/index",
				genMenu3("302", "301", "更新配置", "app:config:edit"),
			),
			// version 版本
			genMenu2("310", "300", "版本管理", "el-icon-Download", 0,
				"app:version:list", "app/version", "app/version/index",
				genMenu3("311", "310", "版本详情", "app:version:detail"),
				genMenu3("312", "310", "版本新增", "app:version:add"),
				genMenu3("313", "310", "版本编辑", "app:version:edit"),
				genMenu3("314", "310", "版本删除", "app:version:del"),
			),
			// shieldword 屏蔽词
			genMenu2("320", "300", "屏蔽词管理", "el-icon-WarningFilled", 0,
				"app:shieldword:list", "app/shieldword", "app/shieldword/index",
				genMenu3("321", "320", "屏蔽词详情", "app:shieldword:detail"),
				genMenu3("322", "320", "屏蔽词新增", "app:shieldword:add"),
				genMenu3("323", "320", "屏蔽词编辑", "app:shieldword:edit"),
				genMenu3("324", "320", "屏蔽词删除", "app:shieldword:del"),
			),
			// vpn VPN
			genMenu2("330", "300", "VPN管理", "el-icon-Magnet", 0,
				"app:vpn:list", "app/vpn", "app/vpn/index",
				genMenu3("331", "330", "VPN详情", "app:vpn:detail"),
				genMenu3("332", "330", "VPN新增", "app:vpn:add"),
				genMenu3("333", "330", "VPN编辑", "app:vpn:edit"),
				genMenu3("334", "330", "VPN删除", "app:vpn:del"),
			),
			// notice 公告
			genMenu2("340", "300", "公告管理", "local-icon-tongzhi", 0,
				"app:notice:list", "app/notice", "app/notice/index",
				genMenu3("341", "340", "公告详情", "app:notice:detail"),
				genMenu3("342", "340", "公告新增", "app:notice:add"),
				genMenu3("343", "340", "公告编辑", "app:notice:edit"),
				genMenu3("344", "340", "公告删除", "app:notice:del"),
			),
			// connection 连接
			genMenu2("350", "300", "连接管理", "el-icon-Connection", 0,
				"app:connection:list", "app/connection", "app/connection/index",
				genMenu3("351", "350", "踢出连接", "app:connection:del"),
			),
			// emoji 表情
			genMenu2("360", "300", "表情管理", "local-icon-biaoqing", 0,
				"app:emoji:list", "app/emoji", "app/emoji/index",
				genMenu3("361", "360", "表情详情", "app:emoji:detail"),
				genMenu3("362", "360", "表情新增", "app:emoji:add"),
				genMenu3("363", "360", "表情编辑", "app:emoji:edit"),
				genMenu3("364", "360", "表情删除", "app:emoji:del"),
			),
			// emojigroup 表情组
			genMenu2("370", "300", "表情组管理", "local-icon-chuangyiwuliao", 0,
				"app:emojigroup:list", "app/emojigroup", "app/emojigroup/index",
				genMenu3("371", "370", "表情组详情", "app:emojigroup:detail"),
				genMenu3("372", "370", "表情组新增", "app:emojigroup:add"),
				genMenu3("373", "370", "表情组编辑", "app:emojigroup:edit"),
				genMenu3("374", "370", "表情组删除", "app:emojigroup:del"),
			),
			// link 外链
			genMenu2("380", "300", "外链管理", "el-icon-Link", 0,
				"app:link:list", "app/link", "app/link/index",
				genMenu3("381", "380", "外链详情", "app:link:detail"),
				genMenu3("382", "380", "外链新增", "app:link:add"),
				genMenu3("383", "380", "外链编辑", "app:link:edit"),
				genMenu3("384", "380", "外链删除", "app:link:del"),
			),
		),
		genMenu1("400", "用户管理", "el-icon-UserFilled", 48, "user",
			// defaultconv 默认会话
			genMenu2("410", "400", "默认会话管理", "el-icon-ChatLineRound", 0,
				"user:defaultconv:list", "user/defaultconv", "user/defaultconv/index",
				genMenu3("411", "410", "默认会话详情", "user:defaultconv:detail"),
				genMenu3("412", "410", "默认会话新增", "user:defaultconv:add"),
				genMenu3("413", "410", "默认会话编辑", "user:defaultconv:edit"),
				genMenu3("414", "410", "默认会话删除", "user:defaultconv:del"),
			),
			// invitationcode 邀请码
			genMenu2("420", "400", "邀请码管理", "el-icon-ElementPlus", 0,
				"user:invitationcode:list", "user/invitationcode", "user/invitationcode/index",
				genMenu3("421", "420", "邀请码详情", "user:invitationcode:detail"),
				genMenu3("422", "420", "邀请码新增", "user:invitationcode:add"),
				genMenu3("423", "420", "邀请码编辑", "user:invitationcode:edit"),
				genMenu3("424", "420", "邀请码删除", "user:invitationcode:del"),
			),
			// ipblacklist IP黑名单
			genMenu2("430", "400", "IP黑名单管理", "el-icon-LocationFilled", 0,
				"user:ipblacklist:list", "user/ipblacklist", "user/ipblacklist/index",
				genMenu3("431", "430", "IP黑名单详情", "user:ipblacklist:detail"),
				genMenu3("432", "430", "IP黑名单新增", "user:ipblacklist:add"),
				genMenu3("433", "430", "IP黑名单编辑", "user:ipblacklist:edit"),
				genMenu3("434", "430", "IP黑名单删除", "user:ipblacklist:del"),
			),
			// ipwhitelist IP白名单
			genMenu2("440", "400", "IP白名单管理", "el-icon-Location", 0,
				"user:ipwhitelist:list", "user/ipwhitelist", "user/ipwhitelist/index",
				genMenu3("441", "440", "IP白名单详情", "user:ipwhitelist:detail"),
				genMenu3("442", "440", "IP白名单新增", "user:ipwhitelist:add"),
				genMenu3("443", "440", "IP白名单编辑", "user:ipwhitelist:edit"),
				genMenu3("444", "440", "IP白名单删除", "user:ipwhitelist:del"),
			),
			// model 模型
			genMenu2("450", "400", "用户管理", "el-icon-UserFilled", 0,
				"user:model:list", "user/model", "user/model/index",
				genMenu3("451", "450", "用户详情", "user:model:detail"),
				genMenu3("452", "450", "用户新增", "user:model:add"),
				genMenu3("453", "450", "用户编辑", "user:model:edit"),
				genMenu3("454", "450", "用户删除", "user:model:del"),
			),
			// loginrecord 登录记录
			genMenu2("460", "400", "登录记录管理", "local-icon-xiangji", 0,
				"user:loginrecord:list", "user/loginrecord", "user/loginrecord/index"),
		),
		genMenu1("500", "群聊管理", "local-icon-weixin", 48, "group",
			// model 模型
			genMenu2("510", "500", "群聊管理", "local-icon-kezizhongxin", 0,
				"group:model:list", "group/model", "group/model/index",
				genMenu3("511", "510", "群聊详情", "group:model:detail"),
				genMenu3("512", "510", "群聊新增", "group:model:add"),
				genMenu3("513", "510", "群聊编辑", "group:model:edit"),
				genMenu3("514", "510", "群聊删除", "group:model:del"),
			),
		),
	})
	for _, menu := range menus {
		insertIfNotFound(tx, menu.Id, menu)
	}
}
