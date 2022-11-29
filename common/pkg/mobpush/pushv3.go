package mobpush

const (
	PUSH_PUSH_URI          = "/v3/push/createPush"
	PUSH_GET_BY_WORKID_URI = "/v3/push/getByWorkId"
	PUSH_GET_BY_WORKNO_URI = "/v3/push/getByWorkno"
	PUSH_CANCEL_TASK_URI   = "/push/drop"
	PUSH_REPLACE_TASK_URI  = "/push/replace"
	PUSH_RECALL_TASK_URI   = "/push/recall"
	PUSH_MULTI_URI         = "/v3/push/createMulti"
)

func (client *PushClient) Push(push Push) ([]byte, error) {
	return GetHTTPClient().PostJSON(client, BASE_URL+PUSH_PUSH_URI, push)
}

func (client *PushClient) PushAll(workNo, title, content string) ([]byte, error) {
	push := NewPushModel(client.AppKey)
	push.setWorkno(workNo)
	push.setTitle(title).setContent(content).setTarget(TARGET_ALL)
	return GetHTTPClient().PostJSON(client, BASE_URL+PUSH_PUSH_URI, push)
}

func (client *PushClient) PushByAlias(workNo, title, content string, alias []string) ([]byte, error) {
	push := NewPushModel(client.AppKey)
	push.setWorkno(workNo)
	push.setTitle(title).setContent(content).setTarget(TARGET_ALIAS)
	push.setAlias(alias)
	return GetHTTPClient().PostJSON(client, BASE_URL+PUSH_PUSH_URI, push)
}

func (client *PushClient) PushByTags(workNo, title, content string, tags []string) ([]byte, error) {
	push := NewPushModel(client.AppKey)
	push.setWorkno(workNo)
	push.setTitle(title).setContent(content).setTarget(TARGET_TAGS)
	push.setTags(tags)
	return GetHTTPClient().PostJSON(client, BASE_URL+PUSH_PUSH_URI, push)
}

func (client *PushClient) PushByRids(workNo, title, content string, rids []string) ([]byte, error) {
	push := NewPushModel(client.AppKey)
	push.setWorkno(workNo)
	push.setTitle(title).setContent(content).setTarget(TARGET_RIDS)
	push.setRids(rids)
	return GetHTTPClient().PostJSON(client, BASE_URL+PUSH_PUSH_URI, push)
}

func (client *PushClient) PushByAreas(workNo, title, content string, areas PushAreas) ([]byte, error) {
	push := NewPushModel(client.AppKey)
	push.setWorkno(workNo)
	push.setTitle(title).setContent(content).setTarget(TARGET_AREAS)
	push.setPushAreas(areas)
	return GetHTTPClient().PostJSON(client, BASE_URL+PUSH_PUSH_URI, push)
}

func (client *PushClient) CancelPushTask(workId string) ([]byte, error) {
	params := client.NewRequestData()
	params["batchId"] = workId
	return GetHTTPClient().PostJSON(client, BASE_URL+PUSH_CANCEL_TASK_URI, params)
}

func (client *PushClient) ReplacePushTask(workId, content string) ([]byte, error) {
	params := client.NewRequestData()
	params["batchId"] = workId
	params["content"] = content
	return GetHTTPClient().PostJSON(client, BASE_URL+PUSH_REPLACE_TASK_URI, params)
}

func (client *PushClient) RecallPushTask(workId string) ([]byte, error) {
	params := client.NewRequestData()
	params["batchId"] = workId
	return GetHTTPClient().PostJSON(client, BASE_URL+PUSH_RECALL_TASK_URI, params)
}

func (client *PushClient) GetPushByBatchId(batchId string) ([]byte, error) {
	params := client.NewRequestData()
	params["workId"] = batchId
	return GetHTTPClient().PostJSON(client, BASE_URL+PUSH_GET_BY_WORKID_URI, params)
}

func (client *PushClient) GetPushByWorkno(workNo string) ([]byte, error) {
	params := client.NewRequestData()
	params["workno"] = workNo
	return GetHTTPClient().PostJSON(client, BASE_URL+PUSH_GET_BY_WORKNO_URI, params)
}

func (client *PushClient) PushMulti(pushMulti PushMulti) ([]byte, error) {
	return GetHTTPClient().PostJSON(client, BASE_URL+PUSH_MULTI_URI, pushMulti)
}
