package xconf

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"strings"
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

func (m *ConfigMgr) initData() {
	// 群聊
	// default_group_description 默认群描述
	m.insertIfNotFound("default_group_description", str("群聊", "default_group_description", "欢迎加入本群", "默认群描述"))
	// default_group_max_member 默认群最大人数
	m.insertIfNotFound("default_group_max_member", num("群聊", "default_group_max_member", 200000, "默认群最大人数"))
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

	// 每个平台ip黑白名单
	ipListModeOption := utils.AnyToString([]map[string]interface{}{
		{
			"label": "禁用",
			"value": "0",
		},
		{
			"label": "黑名单",
			"value": "1",
		},
		{
			"label": "白名单",
			"value": "2",
		},
	})
	// ip_list_mode.ios ios ip黑白名单模式
	for _, platform := range []string{"ios", "android", "windows", "macos", "ipad", "linux", "web"} {
		m.insertIfNotFound("ip_list_mode."+platform, &appmgmtmodel.Config{
			Group:          "用户ip名单模式",
			K:              "ip_list_mode." + platform,
			V:              "0",
			Type:           "option",
			Name:           platform,
			ScopePlatforms: platform,
			Options:        ipListModeOption,
		})
	}
}

func (m *ConfigMgr) DefaultGroupDescription(ctx context.Context) string {
	return m.GetCtx(ctx, "default_group_description")
}

func (m *ConfigMgr) DefaultGroupMaxMember(ctx context.Context) int64 {
	return utils.AnyToInt64(m.GetCtx(ctx, "default_group_max_member"))
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
