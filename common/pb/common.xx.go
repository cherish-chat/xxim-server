package pb

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"
)

func (x *CommonResp) SetMsg(msg string) {
	x.Msg = &msg
}

//func (x *CommonResp) SetTitleMsg(title string, msg string) {
//	buf, _ := json.Marshal(map[string]string{
//		"title": title,
//		"msg":   msg,
//	})
//	msg = string(buf)
//	x.Msg = &msg
//}

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
func NewCommonRequestResp(tip ...string) *CommonResp {
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

var defaultAlertAction = &AlertAction{
	Action: AlertAction_Cancel,
	Title:  "确定",
	JumpTo: "",
}

func NewAlertErrorResp(title string, alert string, actions ...*AlertAction) *CommonResp {
	x := NewCommonResp(CommonResp_AlertError)
	//x.SetTitleMsg(title, alert)
	if len(actions) == 0 {
		actions = append(actions, defaultAlertAction)
	}
	buf, _ := json.Marshal(map[string]any{
		"title":   title,
		"msg":     alert,
		"actions": actions,
	})
	msg := string(buf)
	x.Msg = &msg
	return x
}

func NewRetryErrorResp() *CommonResp {
	return NewCommonResp(CommonResp_RetryError)
}

func (x *CommonResp) Failed() bool {
	return x.Code != CommonResp_Success
}

// Scan 实现 sql 接口
func (x *CommonReq) Scan(input interface{}) error {
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

func PageIndex(pageIndex string) int64 {
	if pageIndex == "" {
		return math.MaxInt64
	}
	parseInt, err := strconv.ParseInt(pageIndex, 10, 64)
	if err != nil {
		return math.MaxInt64
	}
	return parseInt
}

const (
	Platform_Android = "android"
	Platform_IOS     = "ios"
	Platfrom_Ipad    = "ipad"
	Platform_Windows = "windows"
	Platform_Mac     = "macos"
	Platform_Linux   = "linux"
	Platfrom_Web     = "web"
)

var PlatformMap = map[string]string{
	Platform_Android: "安卓",
	Platform_IOS:     "苹果",
	Platfrom_Ipad:    "苹果平板",
	Platform_Windows: "Windows",
	Platform_Mac:     "MacOS",
	Platform_Linux:   "Linux",
	Platfrom_Web:     "Web",
}

func (x *CommonReq) IsAndroid() bool {
	if strings.ToLower(x.Platform) == Platform_Android {
		return true
	}
	return false
}

func (x *CommonReq) IsIOS() bool {
	if strings.ToLower(x.Platform) == Platform_IOS {
		return true
	}
	return false
}

func (x *CommonReq) IsWindows() bool {
	if strings.ToLower(x.Platform) == Platform_Windows {
		return true
	}
	return false
}

func (x *CommonReq) IsMac() bool {
	if strings.ToLower(x.Platform) == Platform_Mac {
		return true
	}
	return false
}

func (x *CommonReq) IsLinux() bool {
	if strings.ToLower(x.Platform) == Platform_Linux {
		return true
	}
	return false
}

func (x *CommonReq) IsWeb() bool {
	if strings.ToLower(x.Platform) == Platfrom_Web {
		return true
	}
	return false
}

func (x *CommonReq) IsMobile() bool {
	return x.IsAndroid() || x.IsIOS()
}

func (x *CommonReq) IsPC() bool {
	return x.IsWindows() || x.IsMac() || x.IsLinux()
}

func (x *CommonReq) IsApp() bool {
	return x.IsMobile() || x.IsPC()
}
