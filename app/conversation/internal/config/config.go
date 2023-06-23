package config

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	RedisConf       redis.RedisConf
	MongoCollection struct {
		Group              xmgo.MongoCollectionConf
		GroupSubscribe     xmgo.MongoCollectionConf
		ConversationMember xmgo.MongoCollectionConf
		Friend             xmgo.MongoCollectionConf
		FriendApplyRecord  xmgo.MongoCollectionConf
		Subscription       xmgo.MongoCollectionConf
	}
	Group struct {
		MinGroupId int `json:",default=100000"`
		//创建时必填名称
		RequiredName bool `json:",default=true"`
		//创建时必填头像
		RequiredAvatar bool `json:",default=true"`
		//默认名称规则
		DefaultNameRule string `json:",options=byMember|fixed"` // byMember:根据成员名称拼接，fixed:固定名称
		//FixedName 固定名称
		FixedName string `json:",default=未命名群聊"` // 固定名称 只有DefaultNameRule=fixed时有效
		//默认头像规则
		DefaultAvatarRule string `json:",options=byName|fixed"` // byName:根据名称生成，fixed:固定头像
		//FixedAvatar 固定头像
		FixedAvatar string `json:",default=group_avatar.png"` // 固定头像 只有DefaultAvatarRule=fixed时有效
		//ByNameAvatarBgColors 根据昵称生成头像的背景颜色 ex: ["#ffffff","#000000"]
		ByNameAvatarBgColors []string
		//ByNameAvatarFgColors 根据昵称生成头像的字体颜色 ex: ["#ffffff","#000000"]
		ByNameAvatarFgColors []string
		//一个人能加入的最大群组数量
		JoinedMaxCount int `json:",default=500"`
		//MaxMemberCount 最大成员数
		MaxMemberCount int `json:",default=200000"` // 最大成员数 默认200000
		Create         struct {
			//AllowPlatform 可以接受的platform创建
			//const (
			//	Platform_IOS        Platform = 0 // ios
			//	Platform_ANDROID    Platform = 1 // android
			//	Platform_WEB        Platform = 2 // web
			//	Platform_WINDOWS    Platform = 3 // windows
			//	Platform_MAC        Platform = 4 // mac
			//	Platform_LINUX      Platform = 5 // linux
			//	Platform_Ipad       Platform = 6 // ipad
			//	Platform_AndroidPad Platform = 7 // android pad
			//)
			AllowPlatform []pb.Platform `json:",optional"`
			//AllowRole 允许什么角色创建
			//const (
			//  AccountRoleUser  = "user"
			//  AccountRoleRobot = "robot"
			//)
			AllowRole []string `json:",optional"`
		}
		Invite struct {
			//AllowPlatform 可以接受的platform邀请
			//const (
			//	Platform_IOS        Platform = 0 // ios
			//	Platform_ANDROID    Platform = 1 // android
			//	Platform_WEB        Platform = 2 // web
			//	Platform_WINDOWS    Platform = 3 // windows
			//	Platform_MAC        Platform = 4 // mac
			//	Platform_LINUX      Platform = 5 // linux
			//	Platform_Ipad       Platform = 6 // ipad
			//	Platform_AndroidPad Platform = 7 // android pad
			//)
			AllowPlatform []pb.Platform `json:",optional"`
			//AllowRole 允许什么角色邀请
			//const (
			//  AccountRoleUser  = "user"
			//  AccountRoleRobot = "robot"
			//)
			AllowRole []string `json:",optional"`

			UserDefaultAllow bool `json:",default=true"` // 用户默认允许邀请
			//DefaultWelcomeMessage 默认欢迎消息
			DefaultWelcomeMessage string `json:",default=欢迎加入群聊"`
		}
		Apply struct {
			//AllowPlatform 可以接受的platform邀请
			//const (
			//	Platform_IOS        Platform = 0 // ios
			//	Platform_ANDROID    Platform = 1 // android
			//	Platform_WEB        Platform = 2 // web
			//	Platform_WINDOWS    Platform = 3 // windows
			//	Platform_MAC        Platform = 4 // mac
			//	Platform_LINUX      Platform = 5 // linux
			//	Platform_Ipad       Platform = 6 // ipad
			//	Platform_AndroidPad Platform = 7 // android pad
			//)
			AllowPlatform []pb.Platform `json:",optional"`
			//AllowRole 允许什么角色邀请
			//const (
			//  AccountRoleUser  = "user"
			//  AccountRoleRobot = "robot"
			//)
			AllowRole []string `json:",optional"`
		}
	}
	RpcClientConf struct {
		Dispatch     zrpc.RpcClientConf
		User         zrpc.RpcClientConf
		Conversation zrpc.RpcClientConf
		Third        zrpc.RpcClientConf
		Message      zrpc.RpcClientConf
	}
	Friend struct {
		DefaultApplySetting string `json:",default={}"`
		//AllowPlatform 可以接受的platform
		//const (
		//	Platform_IOS        Platform = 0 // ios
		//	Platform_ANDROID    Platform = 1 // android
		//	Platform_WEB        Platform = 2 // web
		//	Platform_WINDOWS    Platform = 3 // windows
		//	Platform_MAC        Platform = 4 // mac
		//	Platform_LINUX      Platform = 5 // linux
		//	Platform_Ipad       Platform = 6 // ipad
		//	Platform_AndroidPad Platform = 7 // android pad
		//)
		AllowPlatform []pb.Platform `json:",optional"`
		//AllowRoleApply 允许什么角色申请
		//const (
		//  AccountRoleUser  = "user"
		//  AccountRoleRobot = "robot"
		//)
		AllowRoleApply []string `json:",optional"`
		//AllowRoleBeApplied 允许什么角色被申请
		//const (
		//  AccountRoleUser  = "user"
		//  AccountRoleRobot = "robot"
		//)
		AllowRoleBeApplied []string `json:",optional"`
		//MaxFriendCount 最大好友数
		MaxFriendCount int64 `json:",default=2000"`
		//DefaultSayHello 默认打招呼内容
		DefaultSayHello string `json:",default=我通过了你的好友验证请求，现在我们可以开始聊天了"`
	}
}
