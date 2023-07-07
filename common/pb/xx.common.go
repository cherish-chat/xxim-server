package pb

import (
	"strings"
)

func PlatformFromString(s string) Platform {
	switch strings.ToLower(s) {
	case "ios":
		return Platform_IOS
	case "android":
		return Platform_ANDROID
	case "web":
		return Platform_WEB
	case "windows":
		return Platform_WINDOWS
	case "mac":
		return Platform_MAC
	case "linux":
		return Platform_LINUX
	case "ipad":
		return Platform_Ipad
	case "androidpad":
		return Platform_AndroidPad
	}
	return Platform_WEB
}

func (x Platform) ToString() string {
	switch x {
	case Platform_IOS:
		return "ios"
	case Platform_ANDROID:
		return "android"
	case Platform_WEB:
		return "web"
	case Platform_WINDOWS:
		return "windows"
	case Platform_MAC:
		return "mac"
	case Platform_LINUX:
		return "linux"
	case Platform_Ipad:
		return "ipad"
	case Platform_AndroidPad:
		return "androidpad"
	default:
		return "web"
	}
}
