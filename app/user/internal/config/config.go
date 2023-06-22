package config

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	RedisConf       redis.RedisConf
	MongoCollection struct {
		User        xmgo.MongoCollectionConf
		UserSetting xmgo.MongoCollectionConf
	}
	Account struct {
		//JwtConfig jwt配置
		JwtConfig utils.JwtConfig
		//UsernameUnique 用户名是否唯一
		UsernameUnique bool `json:",optional"`
		//UserRegex 用户名正则
		UserRegex string `json:",optional"`
		//手机号是否唯一
		PhoneUnique bool `json:",optional"`
		//PhoneRegex 手机号正则
		PhoneRegex string `json:",optional"`
		//PhoneCode 国家区号
		PhoneCode []string `json:",optional"`
		//邮箱是否唯一
		EmailUnique bool `json:",optional"`
		//EmailRegex 邮箱正则
		EmailRegex string `json:",optional"`
		Register   struct {
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
			// 默认昵称规则
			DefaultNicknameRule string `json:",options=random|fixed"` // random:随机生成，fixed:固定昵称
			//FixedNickname 固定昵称
			FixedNickname string `json:",default=用户"` // 固定昵称 只有DefaultNicknameRule=fixed时有效
			//RandomNicknamePrefix 随机昵称前缀
			RandomNicknamePrefix string `json:",default=用户"` // 随机昵称前缀 只有DefaultNicknameRule=random时有效
			// 是否必填avatar
			RequireAvatar bool `json:",optional"`
			//默认头像规则
			DefaultAvatarRule string `json:",options=byName|fixed"` // byName:根据昵称生成，fixed:固定头像
			//ByNameAvatarBgColors 根据昵称生成头像的背景颜色 ex: ["#ffffff","#000000"]
			ByNameAvatarBgColors []string
			//ByNameAvatarFgColors 根据昵称生成头像的字体颜色 ex: ["#ffffff","#000000"]
			ByNameAvatarFgColors []string
			//FixedAvatar 固定头像
			FixedAvatar string `json:",default=group_avatar.png"` // 固定头像 只有DefaultAvatarRule=fixed时有效
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
			// RequireCaptcha 是否验证图形验证码
			RequireCaptcha bool `json:",optional"`
		}
		Robot struct {
			// 是否允许用户创建机器人
			AllowCreate bool `json:",optional"`
			// 是否必填nickname
			RequireNickname bool `json:",optional"`
			// DefaultNickname
			DefaultNickname string `json:",default=Robot"`
			// RequireAvatar
			RequireAvatar bool `json:",optional"`
		}
	}
	RpcClientConf struct {
		Dispatch     zrpc.RpcClientConf
		User         zrpc.RpcClientConf
		Conversation zrpc.RpcClientConf
		Third        zrpc.RpcClientConf
		Message      zrpc.RpcClientConf
	}
}
