package mobpush

type PushOperator struct {
	DropType int    `json:"dropType"`
	DropId   string `json:"dropId"`
	NotifyId string `json:"notifyId"`
}
