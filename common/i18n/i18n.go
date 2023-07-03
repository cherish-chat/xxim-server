package i18n

import (
	"encoding/json"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
)

func Get(lang pb.I18NLanguage, zh string) (result string) {
	// TODO: 从配置文件中读取
	return zh
}

func NewServerError(reqHeader *pb.RequestHeader, env string, err error) *pb.ResponseHeader {
	language := pb.I18NLanguage_Chinese_Simplified
	if reqHeader != nil {
		language = reqHeader.GetLanguage()
	}
	data := &pb.ToastActionData{
		Level:   pb.ToastActionData_ERROR,
		Message: Get(language, "服务繁忙，请稍后再试"),
	}
	if env != "pro" && err != nil {
		data = &pb.ToastActionData{
			Level:   pb.ToastActionData_ERROR,
			Message: err.Error(),
		}
	}
	return &pb.ResponseHeader{
		Code:       pb.ResponseCode_SERVER_ERROR,
		ActionType: pb.ResponseActionType_TOAST_ACTION,
		ActionData: utils.Json.MarshalToString(data),
	}
}

func NewInvalidDataError(msg string) *pb.ResponseHeader {
	data := &pb.ToastActionData{
		Level:   pb.ToastActionData_ERROR,
		Message: Get(pb.I18NLanguage_Chinese_Simplified, "参数错误: "+msg),
	}
	return &pb.ResponseHeader{
		Code:       pb.ResponseCode_INVALID_DATA,
		ActionType: pb.ResponseActionType_TOAST_ACTION,
		ActionData: utils.Json.MarshalToString(data),
	}
}

func NewOkHeader() *pb.ResponseHeader {
	return &pb.ResponseHeader{
		Code:       pb.ResponseCode_SUCCESS,
		ActionType: pb.ResponseActionType_NONE_ACTION,
		ActionData: "",
	}
}

func NewToastHeader(level pb.ToastActionData_Level, message string) *pb.ResponseHeader {
	data := &pb.ToastActionData{
		Level:   level,
		Message: message,
	}
	return &pb.ResponseHeader{
		Code:       pb.ResponseCode_INVALID_DATA,
		ActionType: pb.ResponseActionType_TOAST_ACTION,
		ActionData: utils.Json.MarshalToString(data),
		Extra:      "",
	}
}

func NewAuthError(typ pb.AuthErrorType, message string) *pb.ResponseHeader {
	extra := &pb.AuthErrorExtra{
		Type:    typ,
		Message: message,
	}
	buf, _ := json.Marshal(extra)
	return &pb.ResponseHeader{
		Code:  pb.ResponseCode_UNAUTHORIZED,
		Extra: string(buf),
	}
}

func NewForbiddenError() *pb.ResponseHeader {
	return &pb.ResponseHeader{
		Code: pb.ResponseCode_FORBIDDEN,
	}
}
