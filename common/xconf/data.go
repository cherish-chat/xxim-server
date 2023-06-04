package xconf

import (
	"context"
	"encoding/json"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
	"time"
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
	// read_msg_task_interval 插入已读消息 任务 时间间隔
	m.insertIfNotFound("read_msg_task_interval", num("消息", "read_msg_task_interval", 100, "插入已读消息 任务 时间间隔"))
	//read_msg_task_batch_size 插入已读消息 任务 批量大小
	m.insertIfNotFound("read_msg_task_batch_size", num("消息", "read_msg_task_batch_size", 500, "插入已读消息 任务 批量大小"))
	// enable_msg_cleaner 是否开启消息清理器
	m.insertIfNotFound("enable_msg_cleaner", boolean("消息", "enable_msg_cleaner", false, "是否开启消息清理器"))
	// msg_keep_hour 消息保留时间
	m.insertIfNotFound("msg_keep_hour", num("消息", "msg_keep_hour", 24*30, "消息保留时间"))
	// offline.push.allow_disturb 离线推送是否允许打扰
	m.insertIfNotFound("offline.push.allow_disturb", boolean("消息", "offline.push.allow_disturb", false, "离线推送是否允许打扰"))

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
	// sms.send_limit.interval 短信发送间隔
	m.insertIfNotFound("sms.send_limit.interval", num("用户", "sms.send_limit.interval", 60, "短信发送间隔"))
	// sms.send_limit.everyday 每天短信发送次数
	m.insertIfNotFound("sms.send_limit.everyday", num("用户", "sms.send_limit.everyday", 10, "每天短信发送次数"))

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
	// register.mobile.captcha_code 是否使用图形验证码
	m.insertIfNotFound("register.mobile.captcha_code", boolean("注册", "register.mobile.captcha_code", false, "是否使用图形验证码"))
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
	// login.captcha_code 是否使用图形验证码
	m.insertIfNotFound("login.captcha_code", boolean("登录", "login.captcha_code", false, "是否使用图形验证码"))

	// 好友
	// friend.add.user 用户是否能添加用户为好友
	m.insertIfNotFound("friend.add.user", boolean("好友申请", "friend.add.user", true, "用户是否能添加用户为好友"))
	// friend.add.service 用户是否能添加客服为好友
	m.insertIfNotFound("friend.add.service", boolean("好友申请", "friend.add.service", true, "用户是否能添加客服为好友"))
	// friend.add.guest 用户是否能添加游客为好友
	m.insertIfNotFound("friend.add.guest", boolean("好友申请", "friend.add.guest", true, "用户是否能添加游客为好友"))
	// friend.add.robot 用户是否能添加机器人为好友
	m.insertIfNotFound("friend.add.robot", boolean("好友申请", "friend.add.robot", true, "用户是否能添加机器人为好友"))
	// friend_event.include_self 好友事件是否包含自己
	m.insertIfNotFound("friend_event.include_self", boolean("好友申请", "friend_event.include_self", false, "好友事件是否包含自己"))

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
	// message.translate.options 翻译选项
	translateOption := utils.AnyToString([]map[string]interface{}{
		{
			"label": "禁用",
			"value": "0",
		},
		{
			"label": "万维易源",
			"value": "1",
		},
		{
			"label": "Amazon",
			"value": "2",
		},
	})
	m.upsert("message.translate.options", &appmgmtmodel.Config{
		Group:          "文本消息",
		K:              "message.translate.options",
		V:              "0",
		Type:           "option",
		Name:           "翻译选项",
		ScopePlatforms: "",
		Options:        translateOption,
	})
	m.insertIfNotFound("message.translate.languages",
		str("文本消息", "message.translate.languages", "en,hi,id,ja,ko,pt,ru,th,vi,ur,zh,zh-TW", "翻译语言(逗号隔开)"),
	)

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
	// 外部链接地址
	m.insertIfNotFound("group.show_group_info.user.link_url", str("外部链接", "group.show_group_info.user.link_url", "https://www.baidu.com", "外部链接地址"))

	// 发现底部导航按钮名称
	m.insertIfNotFound("group.show_group_info.user.discover_name", str("发现", "group.show_group_info.user.discover_name", "发现", "发现底部导航按钮名称"))
	// 是否显示发现按钮
	m.insertIfNotFound("group.show_group_info.user.discover_show", boolean("发现", "group.show_group_info.user.discover_show", true, "是否显示发现按钮"))
	// 是否允许显示information
	m.insertIfNotFound("information.show.allow", boolean("发现", "information.show.allow", true, "是否允许显示新闻资讯"))

	// 文件上传
	// upload_file_header 上传文件自定义header
	m.insertIfNotFound("upload_file_header", str("文件上传", "upload_file_header", `{}`, "上传文件自定义header"))
	// upload_file_token_secret 上传文件token密钥
	m.insertIfNotFound("upload_file_token_secret", str("文件上传", "upload_file_token_secret", ``, "上传文件token密钥"))
	// upload_file_server_endpoints 上传文件服务器地址
	m.insertIfNotFound("upload_file_server_endpoints", str("文件上传", "upload_file_server_endpoints", `["http://xxx.xxx.xx:80"]`, "上传文件服务器地址"))

	// 万维易源翻译API
	// translate.showapi_appid 万维易源翻译API appid
	m.insertIfNotFound("translate.showapi_appid", str("万维易源翻译API", "translate.showapi_appid", ``, "万维易源翻译API appid"))
	// translate.showapi_sign 万维易源翻译API sign
	m.insertIfNotFound("translate.showapi_sign", str("万维易源翻译API", "translate.showapi_sign", ``, "万维易源翻译API sign"))

	// AmazonTranslate
	// translate.amazon_access_key_id AmazonTranslate access_key_id
	m.insertIfNotFound("translate.amazon_access_key_id", str("AmazonTranslate", "translate.amazon_access_key_id", ``, "AmazonTranslate access_key_id"))
	// translate.amazon_secret_access_key AmazonTranslate secret_access_key
	m.insertIfNotFound("translate.amazon_secret_access_key", str("AmazonTranslate", "translate.amazon_secret_access_key", ``, "AmazonTranslate secret_access_key"))
	// translate.amazon_region
	m.insertIfNotFound("translate.amazon_region", str("AmazonTranslate", "translate.amazon_region", `ap-east-1`, "AmazonTranslate region"))
}

// TranslateOption message.translate.options 翻译选项
func (m *ConfigMgr) TranslateOption(ctx context.Context) string {
	return m.GetCtx(ctx, "message.translate.options", "")
}

// TranslateShowApiAppId translate.showapi_appid 万维易源翻译API appid
func (m *ConfigMgr) TranslateShowApiAppId(ctx context.Context) string {
	return m.GetCtx(ctx, "translate.showapi_appid", "")
}

// TranslateShowApiSign translate.showapi_sign 万维易源翻译API sign
func (m *ConfigMgr) TranslateShowApiSign(ctx context.Context) string {
	return m.GetCtx(ctx, "translate.showapi_sign", "")
}

// TranslateAmazonAccessKeyId translate.amazon_access_key_id AmazonTranslate access_key_id
func (m *ConfigMgr) TranslateAmazonAccessKeyId(ctx context.Context) string {
	return m.GetCtx(ctx, "translate.amazon_access_key_id", "")
}

// TranslateAmazonSecretAccessKey translate.amazon_secret_access_key AmazonTranslate secret_access_key
func (m *ConfigMgr) TranslateAmazonSecretAccessKey(ctx context.Context) string {
	return m.GetCtx(ctx, "translate.amazon_secret_access_key", "")
}

// TranslateAmazonRegion translate.amazon_region AmazonTranslate region
func (m *ConfigMgr) TranslateAmazonRegion(ctx context.Context) string {
	return m.GetCtx(ctx, "translate.amazon_region", "")
}

func (m *ConfigMgr) DefaultGroupDescription(ctx context.Context, userId string) string {
	return m.GetCtx(ctx, "default_group_description", userId)
}

func (m *ConfigMgr) DefaultGroupNewMemberHistoryMsgCount(ctx context.Context, userId string) int64 {
	return utils.AnyToInt64(m.GetCtx(ctx, "default_group_new_member_history_msg_count", userId))
}

func (m *ConfigMgr) DefaultGroupJoinGroupQuestion(ctx context.Context, userId string) string {
	return m.GetCtx(ctx, "default_group_join_group_question", userId)
}

func (m *ConfigMgr) DefaultGroupName(ctx context.Context, userId string) string {
	return m.GetCtx(ctx, "default_group_name", userId)
}

func (m *ConfigMgr) DefaultMinGroupId(ctx context.Context) string {
	return m.GetCtx(ctx, "default_min_group_id", "")
}

func (m *ConfigMgr) DefaultGroupAvatars(ctx context.Context, userId string) []string {
	return m.GetSliceCtx(ctx, "default_group_avatars", userId)
}

func (m *ConfigMgr) OfflinePushTitle(ctx context.Context, userId string) string {
	return m.GetCtx(ctx, "offline_push_title", userId)
}

func (m *ConfigMgr) OfflinePushContent(ctx context.Context, userId string) string {
	return m.GetCtx(ctx, "offline_push_content", userId)
}

// OfflinePushAllowDisturb offline.push.allow_disturb 离线推送是否允许免打扰
func (m *ConfigMgr) OfflinePushAllowDisturb(ctx context.Context, userId string) bool {
	return m.GetCtx(ctx, "offline.push.allow_disturb", userId) == "1"
}

func (m *ConfigMgr) FriendMaxCount(ctx context.Context, userId string) int64 {
	return utils.AnyToInt64(m.GetCtx(ctx, "friend_max_count", userId))
}

func (m *ConfigMgr) NicknameDefault(ctx context.Context) string {
	return m.GetCtx(ctx, "nickname_default", "")
}

func (m *ConfigMgr) SignatureIfNotSet(ctx context.Context) string {
	return m.GetCtx(ctx, "signature_if_not_set", "")
}

func (m *ConfigMgr) AvatarsDefault(ctx context.Context) []string {
	return m.GetSliceCtx(ctx, "avatars_default", "")
}

// UploadFileHeader 上传文件附加header
func (m *ConfigMgr) UploadFileHeader(ctx context.Context) map[string]string {
	v := m.GetCtx(ctx, "upload_file_header", "")
	if v == "" {
		return map[string]string{}
	}
	mp := map[string]string{}
	_ = json.Unmarshal([]byte(v), &mp)
	return mp
}

// UploadFileTokenSecret 上传文件token secret
func (m *ConfigMgr) UploadFileTokenSecret(ctx context.Context) string {
	return m.GetCtx(ctx, "upload_file_token_secret", "")
}

// UploadFileServerEndpoints 上传文件服务器地址
func (m *ConfigMgr) UploadFileServerEndpoints(ctx context.Context) []string {
	v := m.GetCtx(ctx, "upload_file_server_endpoints", "")
	if v == "" {
		return []string{}
	}
	mp := []string{}
	_ = json.Unmarshal([]byte(v), &mp)
	return mp
}

// RegisterIpLimit register.ip_limit.period register.ip_limit.quota
func (m *ConfigMgr) RegisterIpLimit(ctx context.Context) (int, int) {
	period := utils.AnyToInt64(m.GetCtx(ctx, "register.ip_limit.period", ""))
	quota := utils.AnyToInt64(m.GetCtx(ctx, "register.ip_limit.quota", ""))
	return int(period), int(quota)
}

// RegisterMustIpInWhiteList register.must.ip_in_white_list
func (m *ConfigMgr) RegisterMustIpInWhiteList(ctx context.Context) bool {
	return m.GetCtx(ctx, "register.must.ip_in_white_list", "") == "1"
}

// EnablePlatformRegister register.allow.$platform
func (m *ConfigMgr) EnablePlatformRegister(ctx context.Context, platform string) bool {
	if _, ok := pb.PlatformMap[platform]; !ok {
		logx.WithContext(ctx).Errorf("EnablePlatformRegister invalid platform: %s", platform)
		return false
	}
	return m.GetByPlatformCtx(ctx, "register.allow."+platform, platform, "") == "1"
}

// RegisterMustInviteCode register.invite_code.required
func (m *ConfigMgr) RegisterMustInviteCode(ctx context.Context) bool {
	return m.GetCtx(ctx, "register.invite_code.required", "") == "1"
}

// RegisterMustMobile register.mobile.required
func (m *ConfigMgr) RegisterMustMobile(ctx context.Context) bool {
	return m.GetCtx(ctx, "register.mobile.required", "") == "1"
}

// RegisterMustSmsCode register.mobile.sms
func (m *ConfigMgr) RegisterMustSmsCode(ctx context.Context) bool {
	return m.GetCtx(ctx, "register.mobile.sms", "") == "1"
}

// RegisterMustCaptchaCode register.mobile.captcha_code 注册时是否需要填写图形验证码
func (m *ConfigMgr) RegisterMustCaptchaCode(ctx context.Context) bool {
	return m.GetCtx(ctx, "register.mobile.captcha_code", "") == "1"
}

// LoginMustCaptchaCode login.captcha_code 登录时是否需要填写图形验证码
func (m *ConfigMgr) LoginMustCaptchaCode(ctx context.Context) bool {
	return m.GetCtx(ctx, "login.captcha_code", "") == "1"
}

// RegisterMustAvatar register.avatar.required
func (m *ConfigMgr) RegisterMustAvatar(ctx context.Context) bool {
	return m.GetCtx(ctx, "register.avatar.required", "") == "1"
}

// RegisterMustNickname register.nickname.required
func (m *ConfigMgr) RegisterMustNickname(ctx context.Context) bool {
	return m.GetCtx(ctx, "register.nickname.required", "") == "1"
}

// LoginUserOnPlatform login.allow_user.$platform
func (m *ConfigMgr) LoginUserOnPlatform(ctx context.Context, platform string) bool {
	if _, ok := pb.PlatformMap[platform]; !ok {
		logx.WithContext(ctx).Errorf("LoginUserOnPlatform invalid platform: %s", platform)
		return false
	}
	return m.GetByPlatformCtx(ctx, "login.allow_user."+platform, platform, "") == "1"
}

// LoginServiceOnPlatform login.allow_service.$platform
func (m *ConfigMgr) LoginServiceOnPlatform(ctx context.Context, platform string, userId string) bool {
	if _, ok := pb.PlatformMap[platform]; !ok {
		logx.WithContext(ctx).Errorf("LoginUserOnPlatform invalid platform: %s", platform)
		return false
	}
	return m.GetByPlatformCtx(ctx, "login.allow_service."+platform, platform, userId) == "1"
}

// LoginGuestOnPlatform login.allow_guest.$platform
func (m *ConfigMgr) LoginGuestOnPlatform(ctx context.Context, platform string, userId string) bool {
	if _, ok := pb.PlatformMap[platform]; !ok {
		logx.WithContext(ctx).Errorf("LoginUserOnPlatform invalid platform: %s", platform)
		return false
	}
	return m.GetByPlatformCtx(ctx, "login.allow_guest."+platform, platform, userId) == "1"
}

// LoginUserNeedIpWhiteList login.must_user.ip_in_white_list
func (m *ConfigMgr) LoginUserNeedIpWhiteList(ctx context.Context, userId string) bool {
	return m.GetCtx(ctx, "login.must_user.ip_in_white_list", userId) == "1"
}

// LoginServiceNeedIpWhiteList login.must_service.ip_in_white_list
func (m *ConfigMgr) LoginServiceNeedIpWhiteList(ctx context.Context, userId string) bool {
	return m.GetCtx(ctx, "login.must_service.ip_in_white_list", userId) == "1"
}

// LoginGuestNeedIpWhiteList login.must_guest.ip_in_white_list
func (m *ConfigMgr) LoginGuestNeedIpWhiteList(ctx context.Context, userId string) bool {
	return m.GetCtx(ctx, "login.must_guest.ip_in_white_list", userId) == "1"
}

// UserCanAddUserAsFriend friend.add.user
func (m *ConfigMgr) UserCanAddUserAsFriend(ctx context.Context, userId string) bool {
	return m.GetCtx(ctx, "friend.add.user", userId) == "1"
}

// UserCanAddServiceAsFriend friend.add.service
func (m *ConfigMgr) UserCanAddServiceAsFriend(ctx context.Context, userId string) bool {
	return m.GetCtx(ctx, "friend.add.service", userId) == "1"
}

// UserCanAddGuestAsFriend friend.add.guest
func (m *ConfigMgr) UserCanAddGuestAsFriend(ctx context.Context, userId string) bool {
	return m.GetCtx(ctx, "friend.add.guest", userId) == "1"
}

// UserPasswordErrorMaxCount login.password_error_limit
func (m *ConfigMgr) UserPasswordErrorMaxCount(ctx context.Context, userId string) int64 {
	return utils.AnyToInt64(m.GetCtx(ctx, "login.password_error_limit", userId))
}

// SmsErrorMaxCount sms_error_limit.quota
func (m *ConfigMgr) SmsErrorMaxCount(ctx context.Context) int64 {
	return utils.AnyToInt64(m.GetCtx(ctx, "sms_error_limit.quota", ""))
}

// SmsErrorPeriod sms_error_limit.period
func (m *ConfigMgr) SmsErrorPeriod(ctx context.Context) int64 {
	return utils.AnyToInt64(m.GetCtx(ctx, "sms_error_limit.period", ""))
}

// MessageShieldWordCheck message.shield_word.check
func (m *ConfigMgr) MessageShieldWordCheck(ctx context.Context, userId string) bool {
	return m.GetCtx(ctx, "message.shield_word.check", userId) == "1"
}

// MessageShieldWordAllow message.shield_word.allow
func (m *ConfigMgr) MessageShieldWordAllow(ctx context.Context, userId string) bool {
	return m.GetCtx(ctx, "message.shield_word.allow", userId) == "1"
}

// MessageShieldWordAllowReplace message.shield_word.replace
func (m *ConfigMgr) MessageShieldWordAllowReplace(ctx context.Context, userId string) bool {
	return m.GetCtx(ctx, "message.shield_word.replace", userId) == "1"
}

// MessageShieldWordReplace message.shield_word.replace_word
func (m *ConfigMgr) MessageShieldWordReplace(ctx context.Context) rune {
	word := m.GetCtx(ctx, "message.shield_word.replace_word", "")
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
func (m *ConfigMgr) GroupAllowUserQuit(ctx context.Context, userId string) bool {
	return m.GetCtx(ctx, "group.quit_user.allow", userId) == "1"
}

// SmsSendLimitInterval sms.send_limit.interval
func (m *ConfigMgr) SmsSendLimitInterval(ctx context.Context) int {
	return int(utils.AnyToInt64(m.GetCtx(ctx, "sms.send_limit.interval", "")))
}

// SmsSendLimitEveryday sms.send_limit.everyday
func (m *ConfigMgr) SmsSendLimitEveryday(ctx context.Context) int64 {
	return utils.AnyToInt64(m.GetCtx(ctx, "sms.send_limit.everyday", ""))
}

// ReadMsgTaskInterval 已读消息任务间隔
func (m *ConfigMgr) ReadMsgTaskInterval(ctx context.Context) time.Duration {
	return time.Duration(utils.AnyToInt64(m.GetOrDefaultCtx(ctx, "read_msg_task_interval", "100", ""))) * time.Millisecond
}

// ReadMsgTaskBatchSize 已读消息任务批量大小
func (m *ConfigMgr) ReadMsgTaskBatchSize(ctx context.Context) int {
	return int(utils.AnyToInt64(m.GetOrDefaultCtx(ctx, "read_msg_task_batch_size", "500", "")))
}

// EnableMsgCleaner 是否开启消息清理
func (m *ConfigMgr) EnableMsgCleaner() bool {
	return m.GetCtx(context.Background(), "enable_msg_cleaner", "") == "1"
}

// GetMsgKeepHour 消息保留时间
func (m *ConfigMgr) GetMsgKeepHour() int64 {
	return utils.AnyToInt64(m.GetCtx(context.Background(), "msg_keep_hour", ""))
}

// FriendEventIncludeSelf 好友事件是否包含自己
func (m *ConfigMgr) FriendEventIncludeSelf(ctx context.Context) bool {
	return m.GetCtx(ctx, "friend_event.include_self", "") == "1"
}
