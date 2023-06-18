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
		ConversationMember xmgo.MongoCollectionConf
		Friend             xmgo.MongoCollectionConf
		FriendApplyRecord  xmgo.MongoCollectionConf
	}
	Group struct {
		MinGroupId int `json:",default=100000"`
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
	}
}
