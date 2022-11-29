package mobpush

type PushForward struct {
	NextType       int       `json:"nextType"`
	Url            string    `json:"url"`
	Scheme         string    `json:"scheme"`
	SchemeDataList []PushMap `json:"schemeDataList"`
}
