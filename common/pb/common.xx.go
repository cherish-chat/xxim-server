package pb

import (
	"encoding/json"
)

func (x *CommonResp) SetMsg(msg string) {
	x.Msg = &msg
}

func (x *CommonResp) SetTitleMsg(title string, msg string) {
	buf, _ := json.Marshal(map[string]string{
		"title": title,
		"msg":   msg,
	})
	msg = string(buf)
	x.Msg = &msg
}

func (x *CommonResp) setDefaultMsg() {
	msg := ""
	switch x.Code {
	case CommonResp_Success:
		msg = "ok"
	case CommonResp_InternalError:
		msg = "服务繁忙，请稍后再试"
	case CommonResp_RequestError:
		msg = "请求错误"
	case CommonResp_ToastError:
		msg = "服务繁忙，请稍后再试"
	case CommonResp_AlertError:
		msg = `{"title":"提示","msg":"服务繁忙，请稍后再试"}`
	case CommonResp_RetryError:
		msg = "网络繁忙，请重试"
	case CommonResp_UnknownError:
		msg = "未知错误"
	case CommonResp_AuthError:
		msg = "认证失败，请重新登录"
	}
	x.Msg = &msg
}

func NewCommonResp(code CommonResp_Code) *CommonResp {
	x := &CommonResp{
		Code: code,
	}
	x.setDefaultMsg()
	return x
}
func NewSuccessResp() *CommonResp {
	x := NewCommonResp(CommonResp_Success)
	return x
}
func NewUnknownErrorResp(err ...string) *CommonResp {
	x := NewCommonResp(CommonResp_UnknownError)
	if len(err) > 0 {
		x.SetMsg(err[0])
	}
	return x
}
func NewInternalErrorResp(err ...string) *CommonResp {
	x := NewCommonResp(CommonResp_InternalError)
	if len(err) > 0 {
		x.SetMsg(err[0])
	}
	return x
}
func NewRequestErrorResp(tip ...string) *CommonResp {
	x := NewCommonResp(CommonResp_RequestError)
	if len(tip) > 0 {
		x.SetMsg(tip[0])
	}
	return x
}
func NewAuthErrorResp(msg string) *CommonResp {
	x := NewCommonResp(CommonResp_AuthError)
	x.SetMsg(msg)
	return x
}
func NewToastErrorResp(tip ...string) *CommonResp {
	x := NewCommonResp(CommonResp_ToastError)
	if len(tip) > 0 {
		x.SetMsg(tip[0])
	}
	return x
}
func NewAlertErrorResp(title string, alert string) *CommonResp {
	x := NewCommonResp(CommonResp_AlertError)
	x.SetTitleMsg(title, alert)
	return x
}
func NewRetryErrorResp() *CommonResp {
	return NewCommonResp(CommonResp_RetryError)
}

func (x *CommonResp) Failed() bool {
	return x.Code != CommonResp_Success
}

// Scan 实现 sql 接口
func (x *Requester) Scan(input interface{}) error {
	s := string(input.([]byte))
	err := json.Unmarshal([]byte(s), x)
	if err != nil {

	}
	return nil
}

// Scan 实现 sql 接口
func (x *IpRegion) Scan(input interface{}) error {
	s := string(input.([]byte))
	err := json.Unmarshal([]byte(s), x)
	if err != nil {

	}
	return nil
}
