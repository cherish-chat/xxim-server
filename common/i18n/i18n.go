package i18n

import (
	"encoding/json"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	ServerError    = "server_error"
	ParamError     = "param_error"
	PublicKeyError = "public_key_error"
)

func NewServerError(env string, err error) *pb.ResponseHeader {
	data := &pb.ToastActionData{
		Level:   pb.ToastActionData_ERROR,
		Message: ServerError,
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
	logx.Errorf("invalid data error: %s", msg)
	data := &pb.ToastActionData{
		Level:   pb.ToastActionData_ERROR,
		Message: ParamError,
	}
	return &pb.ResponseHeader{
		Code:       pb.ResponseCode_INVALID_DATA,
		ActionType: pb.ResponseActionType_TOAST_ACTION,
		ActionData: utils.Json.MarshalToString(data),
	}
}

func NewInvalidMethodError() *pb.ResponseHeader {
	return &pb.ResponseHeader{
		Code: pb.ResponseCode_INVALID_METHOD,
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
