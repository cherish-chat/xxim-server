package rediskey

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
