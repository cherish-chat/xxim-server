package rediskey

import "net/url"

func UserToken(uid string) string {
	return "h:user_token:" + uid
}

func UserKey(id string) string {
	return "s:model:user:" + id
}

func UserDeviceMapping(userId string) string {
	return "z:user_device_mapping:" + userId
}

func DeviceUserMapping(deviceId string) string {
	return "z:device_user_mapping:" + deviceId
}

func UserDeviceMappingExpire() int64 {
	return 60 * 60 * 24 * 30
}

func RegisterIpLimitKey(ip string) string {
	return "register_ip_limit:" + ip
}

// UserPasswordErrorCountKey 用户密码输入错误次数
func UserPasswordErrorCountKey(uid string) string {
	return "user_password_error_count:" + uid
}

// LatestTurnConvIdKey 上次轮到的会话id
func LatestTurnConvIdKey(invitationCode string) string {
	return "latest_turn_conv_id:" + invitationCode
}

// SmsCodeKey 短信验证码
func SmsCodeKey(scene string, mobile string) string {
	s := "sms_code:" + scene + ":" + mobile
	// url编码
	s = url.QueryEscape(s)
	return s
}

// SmsCodeErrorKey 短信验证码错误次数
func SmsCodeErrorKey(scene string, mobile string) string {
	s := "sms_code_error:" + scene + ":" + mobile
	// url编码
	s = url.QueryEscape(s)
	return s
}
