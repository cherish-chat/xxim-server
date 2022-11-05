package logic

type MockSubscriber struct {
	UserId                 string
	DeviceId               string
	DeviceMaxSeq           int64
	DeviceLatestUpdateTime int64
	XXX                    string
}
type MockSubscribers struct {
	Subscribers []*MockSubscriber
}

// mock 批量获取会话的订阅者
func mockBatchGetConvSubscribers(convIds []string) map[string]*MockSubscribers {
	ret := make(map[string]*MockSubscribers)
	for _, convId := range convIds {
		ret[convId] = &MockSubscribers{Subscribers: []*MockSubscriber{
			{UserId: "1"}, {UserId: "2"}, {UserId: "3"}, {UserId: "4"}, {UserId: "5"},
			{UserId: "6"}, {UserId: "7"}, {UserId: "8"}, {UserId: "9"}, {UserId: "10"},
		}}
	}
	return ret
}
