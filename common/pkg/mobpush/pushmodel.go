package mobpush

type Push struct {
	WorkNo       string        `json:"workno"` // 自定义任务id（用户自自定义生成且唯一、不能重复）
	Source       string        `json:"source"` // 枚举值 webapi
	Appkey       string        `json:"appkey"`
	PushTarget   PushTarget    `json:"pushTarget"`
	PushNotify   PushNotify    `json:"pushNotify"`
	PushOperator *PushOperator `json:"pushOperator,omitempty"`
	PushForward  *PushForward  `json:"pushForward,omitempty"`
}

func NewPushModel(appKey string) *Push {
	push := &Push{}
	push.Appkey = appKey
	push.getDefaultPushModel()
	return push
}

func (push *Push) getDefaultPushModel() *Push {
	push.Source = "webapi"
	push.PushTarget.TagsType = "1"

	push.PushNotify.Plats = []int{1, 2}
	push.PushNotify.Type = 1
	push.PushNotify.IosProduct = 1
	push.PushNotify.OfflineSeconds = 3600

	//push.PushNotify.AndroidNotify.Warn = "12"
	//push.PushNotify.AndroidNotify.Style = 0
	//push.PushNotify.IosNotify.Sound = "default"
	//push.PushForward.NextType = 0

	return push
}

func (push *Push) setWorkno(workno string) *Push {
	push.WorkNo = workno
	return push
}
