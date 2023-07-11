package i18n

import (
	"encoding/json"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	ServerError         = "server_error"
	ParamError          = "param_error"
	PublicKeyError      = "public_key_error"
	RobotNotAllowCreate = "robot_not_allow_create"
	//"robot_nickname_required"
	RobotNicknameRequired = "robot_nickname_required"
	//"robot_avatar_required"
	RobotAvatarRequired = "robot_avatar_required"
	//robot_id_exist
	RobotIdExist = "robot_id_exist"

	//captcha_error
	CaptchaError = "captcha_error"
	//"platform_not_allow"
	PlatformNotAllow = "platform_not_allow"
	// "password_invalid"
	PasswordInvalid = "password_invalid"
	//"sms_code_invalid"
	SmsCodeInvalid = "sms_code_invalid"
	//"phone_invalid"
	PhoneInvalid = "phone_invalid"
	//"email_invalid"
	EmailInvalid = "email_invalid"
	// "email_code_invalid"
	EmailCodeInvalid = "email_code_invalid"
	//"login_failed"
	LoginFailed = "login_failed"

	//username_required
	UsernameRequired = "username_required"
	//username_format_error
	UsernameFormatError = "username_format_error"
	//username_lock_error
	UsernameLockError = "username_lock_error"
	//username_already_exists
	UsernameAlreadyExists = "username_already_exists"
	//password_salt_required
	PasswordSaltRequired = "password_salt_required"
	//password_required
	PasswordRequired = "password_required"
	//phone_required
	PhoneRequired = "phone_required"
	//phone_code_required
	PhoneCodeRequired = "phone_code_required"
	//phone_format_error
	PhoneFormatError = "phone_format_error"
	//phone_code_error
	PhoneCodeError = "phone_code_error"
	//sms_code_required
	SmsCodeRequired = "sms_code_required"
	//sms_code_error
	SmsCodeError = "sms_code_error"
	//phone_lock_error
	PhoneLockError = "phone_lock_error"
	//phone_already_exists
	PhoneAlreadyExists = "phone_already_exists"
	//email_required
	EmailRequired = "email_required"
	//email_format_error
	EmailFormatError = "email_format_error"
	//email_code_required
	EmailCodeRequired = "email_code_required"
	//email_code_error
	EmailCodeError = "email_code_error"
	//email_lock_error
	EmailLockError = "email_lock_error"
	//email_already_exists
	EmailAlreadyExists = "email_already_exists"
	//nickname_required
	NicknameRequired = "nickname_required"
	//avatar_required
	AvatarRequired = "avatar_required"
	//captcha_required
	CaptchaRequired = "captcha_required"
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
