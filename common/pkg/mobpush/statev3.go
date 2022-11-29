package mobpush

const (
	STATS_GET_BY_WORKID_URI  = "/v3/stats/getByWorkId"
	STATS_GET_BY_WORKIDS_URI = "/v3/stats/getByWorkIds"
	STATS_GET_BY_WORKNO_URI  = "/v3/stats/getByWorkno"
	STATS_GET_BY_HOUR_URI    = "/v3/stats/getByHour"
	STATS_GET_BY_DAY_URI     = "/v3/stats/getByDay"
	STATS_GET_BY_DEVICE_URI  = "/v3/stats/getByDevice"
)

func (client *PushClient) GetStatsByWorkId(workId string) ([]byte, error) {
	params := client.NewRequestData()
	params["workId"] = workId
	return GetHTTPClient().PostJSON(client, BASE_URL+STATS_GET_BY_WORKID_URI, params)
}

func (client *PushClient) GetStatsByWorkIds(workIds []string) ([]byte, error) {
	params := client.NewRequestData()
	params["workIds"] = workIds
	return GetHTTPClient().PostJSON(client, BASE_URL+STATS_GET_BY_WORKIDS_URI, params)
}

func (client *PushClient) GetStatsByWorkno(workno string) ([]byte, error) {
	params := client.NewRequestData()
	params["workno"] = workno
	return GetHTTPClient().PostJSON(client, BASE_URL+STATS_GET_BY_WORKNO_URI, params)
}

func (client *PushClient) GetStatsByHour(hour string) ([]byte, error) {
	params := client.NewRequestData()
	params["hour"] = hour
	return GetHTTPClient().PostJSON(client, BASE_URL+STATS_GET_BY_HOUR_URI, params)
}

func (client *PushClient) GetStatsByDay(day string) ([]byte, error) {
	params := client.NewRequestData()
	params["day"] = day
	return GetHTTPClient().PostJSON(client, BASE_URL+STATS_GET_BY_DAY_URI, params)
}

func (client *PushClient) GetStatsByDevice(workId string, pageIndex, pageSize int) ([]byte, error) {
	params := client.NewRequestData()
	params["workId"] = workId
	params["pageIndex"] = pageIndex
	params["pageSize"] = pageSize
	return GetHTTPClient().PostJSON(client, BASE_URL+STATS_GET_BY_DEVICE_URI, params)
}
