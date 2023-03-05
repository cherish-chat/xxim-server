package rediskey

func ServerConfigKey() string {
	return "s:model:server_config"
}

func AppLineConfigKey() string {
	return "s:model:app_line_config"
}

func MSUserToken(uid string) string {
	return "h:ms_user_token:" + uid
}
