package mobpush

const (
	DEVICE_GET_BY_RID   = "/device-v3/getById"
	DEVICE_GET_DIST     = "/device-v3/distribution"
	DEVICE_QUERY_ALIAS  = "/device-v3/getByAlias"
	DEVICE_UPDATE_ALIAS = "/device-v3/upateByAlias"
	DEVICE_UPDATE_TAGS  = "/device-v3/updateTags"
	DEVICE_QUERY_TAGS   = "/device-v3/queryByTags"
)

func (client *PushClient) GetByRid(registrationId string) ([]byte, error) {
	params := client.NewRequestData()
	return GetHTTPClient().PostJSON(client, BASE_URL+DEVICE_GET_BY_RID+"/"+registrationId, params)
}

func (client *PushClient) GetDeviceDistribution() ([]byte, error) {
	params := client.NewRequestData()
	return GetHTTPClient().PostJSON(client, BASE_URL+DEVICE_GET_DIST, params)
}

func (client *PushClient) QueryByAlias(alias string) ([]byte, error) {
	params := client.NewRequestData()
	return GetHTTPClient().PostJSON(client, BASE_URL+DEVICE_QUERY_ALIAS+"/"+alias, params)
}

func (client *PushClient) UpdateAlias(alias, registrationId string) ([]byte, error) {
	params := client.NewRequestData()
	params["alias"] = alias
	params["registrationId"] = registrationId
	return GetHTTPClient().PostJSON(client, BASE_URL+DEVICE_UPDATE_ALIAS, params)
}

func (client *PushClient) UpdateTags(tags []string, registrationId string, opType int) ([]byte, error) {
	params := client.NewRequestData()
	params["tags"] = tags
	params["registrationId"] = registrationId
	params["opType"] = opType
	return GetHTTPClient().PostJSON(client, BASE_URL+DEVICE_UPDATE_TAGS, params)
}

func (client *PushClient) QueryByTags(tags []string) ([]byte, error) {
	params := client.NewRequestData()
	params["tags"] = tags
	return GetHTTPClient().PostJSON(client, BASE_URL+DEVICE_QUERY_TAGS, params)
}
