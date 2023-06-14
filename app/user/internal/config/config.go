package config

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	RedisConf redis.RedisConf
	User      struct {
		MongoCollection xmgo.MongoCollectionConf
	}
	Account struct {
		Register struct {
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
			// 是否必填password
			RequirePassword bool `json:",optional"`
			// 是否必填nickname
			RequireNickname bool `json:",optional"`
			// 是否必填avatar
			RequireAvatar bool `json:",optional"`
			// 是否必须绑定手机号
			RequireBindPhone bool `json:",optional"`
			// 是否必须绑定邮箱
			RequireBindEmail bool `json:",optional"`
			// 是否验证图形验证码
			RequireCaptcha bool `json:",optional"`
		}
		Login struct {
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
		}
	}
	RpcClientConf struct {
		Dispatch zrpc.RpcClientConf
		Third    zrpc.RpcClientConf
	}
}
