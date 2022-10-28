package pb

func (x *CommonResp) SetMsg(msg string) {
	x.Msg = &msg
}

func (x *CommonResp) SetData(data string) {
	x.Data = &data
}

func (x *CommonResp) setDefaultMsg() {
	msg := ""
	switch x.Code {
	case CommonResp_Success:
		msg = "ok"
	case CommonResp_UnknownError:
		msg = "未知错误"
	case CommonResp_InternalError:
		msg = "服务繁忙，请稍后再试"
	case CommonResp_RequestError:
		msg = "请求错误"
	case CommonResp_AuthError:
		msg = "认证失败，请重新登录"
	case CommonResp_ToastError:
		msg = "服务繁忙，请稍后再试"
	case CommonResp_AlertError:
		msg = "服务繁忙，请稍后再试"
	case CommonResp_RetryError:
		msg = "网络繁忙，请重试"
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
func NewSuccessResp(data ...string) *CommonResp {
	x := NewCommonResp(CommonResp_Success)
	if len(data) > 0 {
		x.SetData(data[0])
	}
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
func NewAuthErrorResp(msgData ...string) *CommonResp {
	x := NewCommonResp(CommonResp_AuthError)
	if len(msgData) > 0 {
		x.SetMsg(msgData[0])
	}
	if len(msgData) > 1 {
		x.SetData(msgData[1])
	}
	return x
}
func NewToastErrorResp(tip ...string) *CommonResp {
	x := NewCommonResp(CommonResp_ToastError)
	if len(tip) > 0 {
		x.SetMsg(tip[0])
	}
	return x
}
func NewAlertErrorResp(alert string) *CommonResp {
	x := NewCommonResp(CommonResp_AlertError)
	x.SetMsg(alert)
	return x
}
func NewRetryErrorResp() *CommonResp {
	return NewCommonResp(CommonResp_RetryError)
}
