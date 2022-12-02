package mobpush

const (
	BASE_URL = "http://api.push.mob.com/"
)

type PushClient struct {
	AppKey    string
	AppSecert string
	BaseUrl   string
}

func NewPushClient(appKey, appSecret string) *PushClient {
	return &PushClient{appKey, appSecret, BASE_URL}
}
