package usermodel

import "github.com/cherish-chat/xxim-proto/peerpb"

func GetJwtUniqueKey(x *peerpb.RequestHeader) string {
	// 1. 单点登录
	// return ""
	// 2. 每个平台只能登录一个，意思是我在手机设备A上登录了，那么在手机设备B上登录时，设备A上的token会失效
	//switch x.Platform {
	//case Platform_IOS, Platform_ANDROID, Platform_Ipad, Platform_AndroidPad:
	//	return "mobile"
	//case Platform_WINDOWS, Platform_MAC, Platform_LINUX:
	//	return "pc"
	//case Platform_WEB:
	//	return "web"
	//default:
	//	return "other"
	//}
	// 3. 一个设备同时登录一次该账号，不能重复登录。
	// return x.InstallId
	// 4. 不限制，即使在同一个设备上，也可以登录多次
	//return uuid.New().String()

	// 这里我选择3
	return x.InstallId
}
