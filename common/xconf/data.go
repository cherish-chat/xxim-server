package xconf

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
	"unicode/utf8"
)

func str(groupName string, k string, v string, desc string, platforms ...string) *appmgmtmodel.Config {
	return &appmgmtmodel.Config{
		Group:          groupName,
		K:              k,
		V:              v,
		Type:           "string",
		Name:           desc,
		ScopePlatforms: utils.If(len(platforms) > 0, strings.Join(platforms, ","), ""),
	}
}
func arrayStr(groupName string, k string, v []string, desc string, platforms ...string) *appmgmtmodel.Config {
	return &appmgmtmodel.Config{
		Group:          groupName,
		K:              k,
		V:              utils.AnyToString(v),
		Type:           "stringArray",
		Name:           desc,
		ScopePlatforms: utils.If(len(platforms) > 0, strings.Join(platforms, ","), ""),
	}
}

func num[T utils.NUM](groupName string, k string, v T, desc string, platforms ...string) *appmgmtmodel.Config {
	return &appmgmtmodel.Config{
		Group:          groupName,
		K:              k,
		V:              utils.AnyToString(v),
		Type:           "number",
		Name:           desc,
		ScopePlatforms: utils.If(len(platforms) > 0, strings.Join(platforms, ","), ""),
	}
}

func boolean(groupName string, k string, v bool, desc string, platform ...string) *appmgmtmodel.Config {
	return &appmgmtmodel.Config{
		Group:          groupName,
		K:              k,
		V:              utils.If(v, "1", "0"),
		Type:           "option",
		Name:           desc,
		ScopePlatforms: utils.If(len(platform) > 0, strings.Join(platform, ","), ""),
		Options: utils.AnyToString([]map[string]interface{}{
			{
				"label": "是",
				"value": "1",
			},
			{
				"label": "否",
				"value": "0",
			},
		}),
	}
}

func (m *ConfigMgr) initData() {
	// 群聊
	// default_group_description 默认群描述
	m.insertIfNotFound("default_group_description", str("群聊", "default_group_description", "欢迎加入本群", "默认群描述"))
	m.delete("default_group_max_member")
	// default_group_new_member_history_msg_count 默认群新成员历史消息数量
	m.insertIfNotFound("default_group_new_member_history_msg_count", num("群聊", "default_group_new_member_history_msg_count", 1000, "默认群新成员历史消息数量"))
	// default_group_join_group_question 默认群加群问题
	m.insertIfNotFound("default_group_join_group_question", str("群聊", "default_group_join_group_question", "欢迎加入本群", "默认群加群问题"))
	// default_group_name 默认群名称
	m.insertIfNotFound("default_group_name", str("群聊", "default_group_name", "未命名群聊", "默认群名称"))
	// default_min_group_id 默认最小群id
	m.insertIfNotFound("default_min_group_id", num("群聊", "default_min_group_id", 100000, "默认最小群id"))
	// default_group_avatars 默认群头像
	m.insertIfNotFound("default_group_avatars", arrayStr("群聊", "default_group_avatars", []string{}, "默认群头像"))

	// 消息
	// offline_push_title 离线推送标题
	m.insertIfNotFound("offline_push_title", str("消息", "offline_push_title", "默认线路", "离线推送标题"))
	// offline_push_content 离线推送内容
	m.insertIfNotFound("offline_push_content", str("消息", "offline_push_content", "您有一条新消息", "离线推送内容"))

	// 好友
	// friend_max_count 好友最大数量
	m.insertIfNotFound("friend_max_count", num("好友", "friend_max_count", 200000, "好友最大数量"))

	// 用户
	// nickname_default 用户默认昵称
	m.insertIfNotFound("nickname_default", str("用户", "nickname_default", "用户", "用户默认昵称"))
	// signature_if_not_set 用户未设置签名时的签名
	m.insertIfNotFound("signature_if_not_set", str("用户", "signature_if_not_set", "这个人很懒，什么都没留下", "用户未设置签名时的签名"))
	// avatars_default 用户默认头像
	m.insertIfNotFound("avatars_default", arrayStr("用户", "avatars_default", []string{}, "用户默认头像"))
	// sms_error_limit.quota 短信错误限制次数
	m.insertIfNotFound("sms_error_limit.quota", num("用户", "sms_error_limit.quota", 5, "短信错误限制次数"))
	// sms_error_limit.period 短信错误限制周期
	m.insertIfNotFound("sms_error_limit.period", num("用户", "sms_error_limit.period", 60, "短信错误限制周期"))

	// 注册
	// register.ip_limit.period 注册ip限制周期 单位秒
	m.insertIfNotFound("register.ip_limit.period", num("注册", "register.ip_limit.period", 60, "注册ip限制周期(秒)"))
	// register.ip_limit.quota 注册ip限制次数
	m.insertIfNotFound("register.ip_limit.quota", num("注册", "register.ip_limit.quota", 10, "注册ip限制次数"))
	// register.allow. 是否允许注册
	registerAllowOption := utils.AnyToString([]map[string]interface{}{
		{
			"label": "允许",
			"value": "1",
		},
		{
			"label": "不允许",
			"value": "0",
		},
	})
	for _, platform := range []string{"ios", "android", "windows", "macos", "ipad", "linux", "web"} {
		m.insertIfNotFound("register.allow."+platform, &appmgmtmodel.Config{
			Group:          "允许注册的平台",
			K:              "register.allow." + platform,
			V:              "0",
			Type:           "option",
			Name:           platform,
			ScopePlatforms: platform,
			Options:        registerAllowOption,
		})
	}
	// register.invite_code.show 是否显示邀请码
	m.insertIfNotFound("register.invite_code.show", boolean("注册", "register.invite_code.show", false, "是否显示邀请码"))
	// register.invite_code.required 是否必填邀请码
	m.insertIfNotFound("register.invite_code.required", boolean("注册", "register.invite_code.required", false, "是否必填邀请码"))
	// register.avatar.show 是否显示头像
	m.insertIfNotFound("register.avatar.show", boolean("注册", "register.avatar.show", false, "是否显示头像"))
	// register.avatar.required 是否必填头像
	m.insertIfNotFound("register.avatar.required", boolean("注册", "register.avatar.required", false, "是否必填头像"))
	// register.nickname.show 是否显示昵称
	m.insertIfNotFound("register.nickname.show", boolean("注册", "register.nickname.show", false, "是否显示昵称"))
	// register.nickname.required 是否必填昵称
	m.insertIfNotFound("register.nickname.required", boolean("注册", "register.nickname.required", false, "是否必填昵称"))
	// register.mobile.required 是否必填手机号
	m.insertIfNotFound("register.mobile.required", boolean("注册", "register.mobile.required", false, "是否必填手机号"))
	// register.mobile.sms 是否使用短信验证码
	m.insertIfNotFound("register.mobile.sms", boolean("注册", "register.mobile.sms", false, "是否使用短信验证码"))
	// register.must.ip_in_white_list ip必须在白名单中
	m.insertIfNotFound("register.must.ip_in_white_list", boolean("注册", "register.must.ip_in_white_list", false, "ip必须在白名单中"))

	// 用户能在哪些平台登录
	// login.allow_user.ios 是否允许ios登录
	for _, platform := range []string{"ios", "android", "windows", "macos", "ipad", "linux", "web"} {
		m.insertIfNotFound("login.allow_user."+platform, &appmgmtmodel.Config{
			Group:          "允许用户登录",
			K:              "login.allow_user." + platform,
			V:              "0",
			Type:           "option",
			Name:           platform,
			ScopePlatforms: platform,
			Options:        registerAllowOption,
		})
	}
	// 客服能在哪些平台登录
	// login.allow_service.ios 是否允许ios登录
	for _, platform := range []string{"ios", "android", "windows", "macos", "ipad", "linux", "web"} {
		m.insertIfNotFound("login.allow_service."+platform, &appmgmtmodel.Config{
			Group:          "允许客服登录",
			K:              "login.allow_service." + platform,
			V:              "0",
			Type:           "option",
			Name:           platform,
			ScopePlatforms: platform,
			Options:        registerAllowOption,
		})
	}
	// 游客能在哪些平台登录
	// login.allow_guest.ios 是否允许ios登录
	for _, platform := range []string{"ios", "android", "windows", "macos", "ipad", "linux", "web"} {
		m.insertIfNotFound("login.allow_guest."+platform, &appmgmtmodel.Config{
			Group:          "允许游客登录",
			K:              "login.allow_guest." + platform,
			V:              "0",
			Type:           "option",
			Name:           platform,
			ScopePlatforms: platform,
			Options:        registerAllowOption,
		})
	}
	// login.must_user.ip_in_white_list 用户登录时ip必须在白名单中
	m.insertIfNotFound("login.must_user.ip_in_white_list", boolean("登录", "login.must_user.ip_in_white_list", false, "用户登录时ip必须在白名单中"))
	// login.must_service.ip_in_white_list 客服登录时ip必须在白名单中
	m.insertIfNotFound("login.must_service.ip_in_white_list", boolean("登录", "login.must_service.ip_in_white_list", false, "客服登录时ip必须在白名单中"))
	// login.must_guest.ip_in_white_list 游客登录时ip必须在白名单中
	m.insertIfNotFound("login.must_guest.ip_in_white_list", boolean("登录", "login.must_guest.ip_in_white_list", false, "游客登录时ip必须在白名单中"))
	// login.password_error_limit 登录密码错误上限
	m.insertIfNotFound("login.password_error_limit", num("登录", "login.password_error_limit", 5, "登录密码错误上限"))

	// 好友
	// friend.add.user 用户是否能添加用户为好友
	m.insertIfNotFound("friend.add.user", boolean("好友申请", "friend.add.user", true, "用户是否能添加用户为好友"))
	// friend.add.service 用户是否能添加客服为好友
	m.insertIfNotFound("friend.add.service", boolean("好友申请", "friend.add.service", true, "用户是否能添加客服为好友"))
	// friend.add.guest 用户是否能添加游客为好友
	m.insertIfNotFound("friend.add.guest", boolean("好友申请", "friend.add.guest", true, "用户是否能添加游客为好友"))
	// friend.add.robot 用户是否能添加机器人为好友
	m.insertIfNotFound("friend.add.robot", boolean("好友申请", "friend.add.robot", true, "用户是否能添加机器人为好友"))

	// 创建群聊
	// group.create.user 用户是否能创建群聊
	m.insertIfNotFound("group.create.user", boolean("群聊", "group.create.user", true, "用户是否能创建群聊"))
	// group.create.service 客服是否能创建群聊
	m.insertIfNotFound("group.create.service", boolean("群聊", "group.create.service", true, "客服是否能创建群聊"))
	// group.create.guest 游客是否能创建群聊
	m.insertIfNotFound("group.create.guest", boolean("群聊", "group.create.guest", true, "游客是否能创建群聊"))

	// 消息页面开关
	// message.show.online_status 是否显示在线状态
	m.insertIfNotFound("message.show.online_status", boolean("消息页面", "message.show.online_status", true, "是否显示在线状态"))
	// message.show.read_status 是否显示已读状态
	m.insertIfNotFound("message.show.read_status", boolean("消息页面", "message.show.read_status", true, "是否显示已读状态"))
	// message.show.typing_status 是否显示正在输入状态
	m.insertIfNotFound("message.show.typing_status", boolean("消息页面", "message.show.typing_status", true, "是否显示正在输入状态"))
	// message.show.chat_log_button 是否显示聊天记录按钮
	m.insertIfNotFound("message.show.chat_log_button", boolean("消息页面", "message.show.chat_log_button", true, "是否显示聊天记录按钮"))
	// message.message_limit.time 消息发送时间间隔
	m.insertIfNotFound("message.message_limit.time", num("消息页面", "message.message_limit.time", 1, "普通消息发送时间间隔"))
	// message.message_time_tip.time 消息时间提示间隔
	m.insertIfNotFound("message.message_time_tip.time", num("消息页面", "message.message_time_tip.time", 300, "消息时间提示间隔"))

	// message.same_message_limit.allow 是否允许发送相同的消息
	m.insertIfNotFound("message.same_message_limit.allow", boolean("消息页面-相同消息", "message.same_message_limit.allow", true, "是否允许发送相同的消息"))
	// message.same_message_limit 相同消息发送上限
	m.insertIfNotFound("message.same_message_limit", num("消息页面-相同消息", "message.same_message_limit", 3, "相同消息发送上限"))
	// message.same_message_limit.time 相同消息发送时间间隔
	m.insertIfNotFound("message.same_message_limit.time", num("消息页面-相同消息", "message.same_message_limit.time", 5, "相同消息发送时间间隔"))

	// message.message_length_limit.text 文本消息长度限制
	m.insertIfNotFound("message.message_length_limit.text", num("文本消息", "message.message_length_limit.text", 500, "文本消息长度限制"))
	// message.shield_word.check 是否开启敏感词检测
	m.insertIfNotFound("message.shield_word.check", boolean("文本消息", "message.shield_word.check", true, "是否开启敏感词检测"))
	// message.shield_word.allow 是否允许发送敏感词
	m.insertIfNotFound("message.shield_word.allow", boolean("文本消息", "message.shield_word.allow", true, "是否允许发送敏感词"))
	// message.shield_word.replace 是否替换敏感词
	m.insertIfNotFound("message.shield_word.replace", boolean("文本消息", "message.shield_word.replace", true, "是否替换敏感词"))
	// message.shield_word.replace_word 替换敏感词
	m.insertIfNotFound("message.shield_word.replace_word", str("文本消息", "message.shield_word.replace_word", "*", "替换敏感词"))

	// 群聊页面
	// group.search.allow 是否允许搜索群聊
	m.insertIfNotFound("group.search.allow", boolean("群聊页面", "group.search.allow", true, "是否允许搜索群聊"))
	// group.quit_user.allow 普通用户是否允许退出群聊
	m.insertIfNotFound("group.quit_user.allow", boolean("群聊页面", "group.quit_user.allow", true, "普通用户是否允许退出群聊"))
	// group.show_member.allow 是否显示群成员
	m.insertIfNotFound("group.show_member.allow", boolean("群聊页面", "group.show_member.allow", true, "是否显示群成员"))
	// group.owner_clear_screen.allow 群主是否允许清屏
	m.insertIfNotFound("group.owner_clear_screen.allow", boolean("群聊页面", "group.owner_clear_screen.allow", true, "群主是否允许清屏"))
	// group.show_member_count.allow 是否显示群成员数量
	m.insertIfNotFound("group.show_member_count.allow", boolean("群聊页面", "group.show_member_count.allow", true, "是否显示群成员数量"))
	// group.show_invite_msg.allow 是否显示邀请入群信息
	m.insertIfNotFound("group.show_invite_msg.allow", boolean("群聊页面", "group.show_invite_msg.allow", true, "是否显示邀请入群信息"))
	// group.show_fake_member_count.allow 是否显示群成员假数量
	m.insertIfNotFound("group.show_fake_member_count.allow", boolean("群聊页面", "group.show_fake_member_count.allow", true, "是否显示群成员假数量"))
	// group.show_group_info.user 普通用户是否显示群信息
	m.insertIfNotFound("group.show_group_info.user", boolean("群聊页面", "group.show_group_info.user", true, "普通用户是否显示群信息"))

	// 外部链接底部导航按钮名称
	m.insertIfNotFound("group.show_group_info.user.link_name", str("外部链接", "group.show_group_info.user.link_name", "发现", "外部链接底部导航按钮名称"))
	// 是否显示外部链接按钮
	m.insertIfNotFound("group.show_group_info.user.link_show", boolean("外部链接", "group.show_group_info.user.link_show", true, "是否显示外部链接按钮"))
}

func (m *ConfigMgr) DefaultGroupDescription(ctx context.Context) string {
	return m.GetCtx(ctx, "default_group_description")
}

func (m *ConfigMgr) DefaultGroupNewMemberHistoryMsgCount(ctx context.Context) int64 {
	return utils.AnyToInt64(m.GetCtx(ctx, "default_group_new_member_history_msg_count"))
}

func (m *ConfigMgr) DefaultGroupJoinGroupQuestion(ctx context.Context) string {
	return m.GetCtx(ctx, "default_group_join_group_question")
}

func (m *ConfigMgr) DefaultGroupName(ctx context.Context) string {
	return m.GetCtx(ctx, "default_group_name")
}

func (m *ConfigMgr) DefaultMinGroupId(ctx context.Context) string {
	return m.GetCtx(ctx, "default_min_group_id")
}

func (m *ConfigMgr) DefaultGroupAvatars(ctx context.Context) []string {
	return m.GetSliceCtx(ctx, "default_group_avatars")
}

func (m *ConfigMgr) OfflinePushTitle(ctx context.Context) string {
	return m.GetCtx(ctx, "offline_push_title")
}

func (m *ConfigMgr) OfflinePushContent(ctx context.Context) string {
	return m.GetCtx(ctx, "offline_push_content")
}

func (m *ConfigMgr) FriendMaxCount(ctx context.Context) int64 {
	return utils.AnyToInt64(m.GetCtx(ctx, "friend_max_count"))
}

func (m *ConfigMgr) NicknameDefault(ctx context.Context) string {
	return m.GetCtx(ctx, "nickname_default")
}

func (m *ConfigMgr) SignatureIfNotSet(ctx context.Context) string {
	return m.GetCtx(ctx, "signature_if_not_set")
}

func (m *ConfigMgr) AvatarsDefault(ctx context.Context) []string {
	return m.GetSliceCtx(ctx, "avatars_default")
}

// RegisterIpLimit register.ip_limit.period register.ip_limit.quota
func (m *ConfigMgr) RegisterIpLimit(ctx context.Context) (int, int) {
	period := utils.AnyToInt64(m.GetCtx(ctx, "register.ip_limit.period"))
	quota := utils.AnyToInt64(m.GetCtx(ctx, "register.ip_limit.quota"))
	return int(period), int(quota)
}

// RegisterMustIpInWhiteList register.must.ip_in_white_list
func (m *ConfigMgr) RegisterMustIpInWhiteList(ctx context.Context) bool {
	return m.GetCtx(ctx, "register.must.ip_in_white_list") == "1"
}

// EnablePlatformRegister register.allow.$platform
func (m *ConfigMgr) EnablePlatformRegister(ctx context.Context, platform string) bool {
	if _, ok := pb.PlatformMap[platform]; !ok {
		logx.WithContext(ctx).Errorf("EnablePlatformRegister invalid platform: %s", platform)
		return false
	}
	return m.GetByPlatformCtx(ctx, "register.allow."+platform, platform) == "1"
}

// RegisterMustInviteCode register.invite_code.required
func (m *ConfigMgr) RegisterMustInviteCode(ctx context.Context) bool {
	return m.GetCtx(ctx, "register.invite_code.required") == "1"
}

// RegisterMustMobile register.mobile.required
func (m *ConfigMgr) RegisterMustMobile(ctx context.Context) bool {
	return m.GetCtx(ctx, "register.mobile.required") == "1"
}

// RegisterMustSmsCode register.mobile.sms
func (m *ConfigMgr) RegisterMustSmsCode(ctx context.Context) bool {
	return m.GetCtx(ctx, "register.mobile.sms") == "1"
}

// RegisterMustAvatar register.avatar.required
func (m *ConfigMgr) RegisterMustAvatar(ctx context.Context) bool {
	return m.GetCtx(ctx, "register.avatar.required") == "1"
}

// RegisterMustNickname register.nickname.required
func (m *ConfigMgr) RegisterMustNickname(ctx context.Context) bool {
	return m.GetCtx(ctx, "register.nickname.required") == "1"
}

// LoginUserOnPlatform login.allow_user.$platform
func (m *ConfigMgr) LoginUserOnPlatform(ctx context.Context, platform string) bool {
	if _, ok := pb.PlatformMap[platform]; !ok {
		logx.WithContext(ctx).Errorf("LoginUserOnPlatform invalid platform: %s", platform)
		return false
	}
	return m.GetByPlatformCtx(ctx, "login.allow_user."+platform, platform) == "1"
}

// LoginServiceOnPlatform login.allow_service.$platform
func (m *ConfigMgr) LoginServiceOnPlatform(ctx context.Context, platform string) bool {
	if _, ok := pb.PlatformMap[platform]; !ok {
		logx.WithContext(ctx).Errorf("LoginUserOnPlatform invalid platform: %s", platform)
		return false
	}
	return m.GetByPlatformCtx(ctx, "login.allow_service."+platform, platform) == "1"
}

// LoginGuestOnPlatform login.allow_guest.$platform
func (m *ConfigMgr) LoginGuestOnPlatform(ctx context.Context, platform string) bool {
	if _, ok := pb.PlatformMap[platform]; !ok {
		logx.WithContext(ctx).Errorf("LoginUserOnPlatform invalid platform: %s", platform)
		return false
	}
	return m.GetByPlatformCtx(ctx, "login.allow_guest."+platform, platform) == "1"
}

// LoginUserNeedIpWhiteList login.must_user.ip_in_white_list
func (m *ConfigMgr) LoginUserNeedIpWhiteList(ctx context.Context) bool {
	return m.GetCtx(ctx, "login.must_user.ip_in_white_list") == "1"
}

// LoginServiceNeedIpWhiteList login.must_service.ip_in_white_list
func (m *ConfigMgr) LoginServiceNeedIpWhiteList(ctx context.Context) bool {
	return m.GetCtx(ctx, "login.must_service.ip_in_white_list") == "1"
}

// LoginGuestNeedIpWhiteList login.must_guest.ip_in_white_list
func (m *ConfigMgr) LoginGuestNeedIpWhiteList(ctx context.Context) bool {
	return m.GetCtx(ctx, "login.must_guest.ip_in_white_list") == "1"
}

// UserCanAddUserAsFriend friend.add.user
func (m *ConfigMgr) UserCanAddUserAsFriend(ctx context.Context) bool {
	return m.GetCtx(ctx, "friend.add.user") == "1"
}

// UserCanAddServiceAsFriend friend.add.service
func (m *ConfigMgr) UserCanAddServiceAsFriend(ctx context.Context) bool {
	return m.GetCtx(ctx, "friend.add.service") == "1"
}

// UserCanAddGuestAsFriend friend.add.guest
func (m *ConfigMgr) UserCanAddGuestAsFriend(ctx context.Context) bool {
	return m.GetCtx(ctx, "friend.add.guest") == "1"
}

// UserPasswordErrorMaxCount login.password_error_limit
func (m *ConfigMgr) UserPasswordErrorMaxCount(ctx context.Context) int64 {
	return utils.AnyToInt64(m.GetCtx(ctx, "login.password_error_limit"))
}

// SmsErrorMaxCount sms_error_limit.quota
func (m *ConfigMgr) SmsErrorMaxCount(ctx context.Context) int64 {
	return utils.AnyToInt64(m.GetCtx(ctx, "sms_error_limit.quota"))
}

// SmsErrorPeriod sms_error_limit.period
func (m *ConfigMgr) SmsErrorPeriod(ctx context.Context) int64 {
	return utils.AnyToInt64(m.GetCtx(ctx, "sms_error_limit.period"))
}

// MessageShieldWordCheck message.shield_word.check
func (m *ConfigMgr) MessageShieldWordCheck(ctx context.Context) bool {
	return m.GetCtx(ctx, "message.shield_word.check") == "1"
}

// MessageShieldWordAllow message.shield_word.allow
func (m *ConfigMgr) MessageShieldWordAllow(ctx context.Context) bool {
	return m.GetCtx(ctx, "message.shield_word.allow") == "1"
}

// MessageShieldWordAllowReplace message.shield_word.replace
func (m *ConfigMgr) MessageShieldWordAllowReplace(ctx context.Context) bool {
	return m.GetCtx(ctx, "message.shield_word.replace") == "1"
}

// MessageShieldWordReplace message.shield_word.replace_word
func (m *ConfigMgr) MessageShieldWordReplace(ctx context.Context) rune {
	word := m.GetCtx(ctx, "message.shield_word.replace_word")
	if len(word) == 0 {
		return '*'
	}
	// 能否转换为rune
	if _, size := utf8.DecodeRuneInString(word); size == 0 {
		return '*'
	} else if size != len(word) {
		return '*'
	}
	return rune(word[0])
}

// GroupAllowUserQuit group.quit_user.allow
func (m *ConfigMgr) GroupAllowUserQuit(ctx context.Context) bool {
	return m.GetCtx(ctx, "group.quit_user.allow") == "1"
}
