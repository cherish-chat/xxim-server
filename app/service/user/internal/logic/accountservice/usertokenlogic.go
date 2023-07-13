package accountservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"go.mongodb.org/mongo-driver/bson"
	"time"

	"github.com/cherish-chat/xxim-server/app/service/user/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserTokenLogic {
	return &UserTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UserToken 用户登录
func (l *UserTokenLogic) UserToken(in *peerpb.UserTokenReq) (*peerpb.UserTokenResp, error) {
	_, passwordExists := in.AccountMap[peerpb.AccountType_Password.String()]
	_, usernameExists := in.AccountMap[peerpb.AccountType_Username.String()]
	_, phoneExists := in.AccountMap[peerpb.AccountType_Phone.String()]
	_, phoneCountryCodeExists := in.AccountMap[peerpb.AccountType_PhoneCountryCode.String()]
	_, emailExists := in.AccountMap[peerpb.AccountType_Email.String()]
	_, smsCodeExists := in.VerifyMap[peerpb.AccountVerifyType_SmsCode.String()]
	_, emailCodeExists := in.VerifyMap[peerpb.AccountVerifyType_EmailCode.String()]

	_, captchaIdExists := in.VerifyMap[peerpb.AccountVerifyType_CaptchaId.String()]
	_, captchaCodeExists := in.VerifyMap[peerpb.AccountVerifyType_CaptchaCode.String()]
	if !captchaIdExists || !captchaCodeExists {
		// 获取token是否需要图形验证码
		if l.svcCtx.Config.Account.Login.RequireCaptcha {
			// 参数错误
			return &peerpb.UserTokenResp{
				Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.ParamError),
			}, nil
		}
	} else if captchaIdExists && captchaCodeExists {
		// 验证图形验证码
		captchaVerifyResp, err := l.svcCtx.CaptchaService.CaptchaVerify(context.Background(), &peerpb.CaptchaVerifyReq{
			Header:      in.Header,
			CaptchaId:   in.VerifyMap[peerpb.AccountVerifyType_CaptchaId.String()],
			CaptchaCode: in.VerifyMap[peerpb.AccountVerifyType_CaptchaCode.String()],
			Delete:      true,
		})
		if err != nil {
			return nil, err
		}
		if !captchaVerifyResp.Success {
			return &peerpb.UserTokenResp{
				Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.CaptchaError),
			}, nil
		}
	}

	//平台
	if !utils.EnumInSlice[peerpb.Platform](in.Header.Platform, l.svcCtx.Config.Account.Login.AllowPlatform) {
		return &peerpb.UserTokenResp{
			Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.PlatformNotAllow),
		}, nil
	}

	if passwordExists {
		// 密码登录，看看是用户名+密码、手机号+密码、还是邮箱+密码
		if usernameExists {
			return l.LoginByPasswordUsername(context.Background(), in)
		}
		if phoneExists && phoneCountryCodeExists {
			return l.LoginByPasswordPhone(context.Background(), in)
		}
		if emailExists {
			return l.LoginByPasswordEmail(context.Background(), in)
		}
		// 参数错误
		return &peerpb.UserTokenResp{
			Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, "参数错误, 密码登录时必须传入用户名、手机号或邮箱"),
		}, nil
	} else {
		// 看看smsCode是否存在，存在则是手机号+验证码登录
		if smsCodeExists {
			// 判断手机号是否存在
			if phoneExists && phoneCountryCodeExists {
				return l.LoginBySmsCode(context.Background(), in)
			}
			// 参数错误
			return &peerpb.UserTokenResp{
				Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, "参数错误, 手机号+验证码登录时必须传入手机号和国家代码"),
			}, nil
		}
		// 看看emailCode是否存在，存在则是邮箱+验证码登录
		if emailCodeExists {
			// 判断邮箱是否存在
			if emailExists {
				return l.LoginByEmailCode(context.Background(), in)
			}
			// 参数错误
			return &peerpb.UserTokenResp{
				Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, "参数错误, 邮箱+验证码登录时必须传入邮箱"),
			}, nil
		}
		// 参数错误
		return &peerpb.UserTokenResp{
			Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, "参数错误, 验证码登录时必须传入验证码"),
		}, nil
	}
}

// LoginByPasswordUsername 用户名密码登录
func (l *UserTokenLogic) LoginByPasswordUsername(ctx context.Context, in *peerpb.UserTokenReq) (*peerpb.UserTokenResp, error) {
	password, username := in.AccountMap[peerpb.AccountType_Password.String()], in.AccountMap[peerpb.AccountType_Username.String()]
	// 通过用户名获取用户信息
	user := &usermodel.User{}
	err := l.svcCtx.UserCollection.Find(context.Background(), bson.M{"accountMap." + peerpb.AccountType_Username.String(): username}).One(user)
	if err != nil {
		l.Errorf("login by password username error: %v", err)
		return &peerpb.UserTokenResp{
			Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.PasswordInvalid),
		}, nil
	}
	// 验证密码
	ok := utils.Pwd.VerifyPwd(password, user.GetAccountMap().Get(peerpb.AccountType_Password.String()), user.GetAccountMap().Get(peerpb.AccountType_PasswordSalt.String()))
	if !ok {
		return &peerpb.UserTokenResp{
			Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.PasswordInvalid),
		}, nil
	}
	// 生成token
	return l.generateToken(in, user), nil
}

// LoginByPasswordPhone 手机号密码登录
func (l *UserTokenLogic) LoginByPasswordPhone(ctx context.Context, in *peerpb.UserTokenReq) (*peerpb.UserTokenResp, error) {
	password, phone, phoneCountryCode := in.AccountMap[peerpb.AccountType_Password.String()], in.AccountMap[peerpb.AccountType_Phone.String()], in.AccountMap[peerpb.AccountType_PhoneCountryCode.String()]
	// 通过手机号获取用户信息
	user := &usermodel.User{}
	err := l.svcCtx.UserCollection.Find(context.Background(), bson.M{"accountMap." + peerpb.AccountType_Phone.String(): phone, "accountMap." + peerpb.AccountType_PhoneCountryCode.String(): phoneCountryCode}).One(user)
	if err != nil {
		l.Errorf("login by password phone error: %v", err)
		return &peerpb.UserTokenResp{
			Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.PasswordInvalid),
		}, nil
	}
	// 验证密码
	ok := utils.Pwd.VerifyPwd(password, user.GetAccountMap().Get(peerpb.AccountType_Password.String()), user.GetAccountMap().Get(peerpb.AccountType_PasswordSalt.String()))
	if !ok {
		return &peerpb.UserTokenResp{
			Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.PasswordInvalid),
		}, nil
	}
	// 生成token
	return l.generateToken(in, user), nil
}

// LoginByPasswordEmail 邮箱密码登录
func (l *UserTokenLogic) LoginByPasswordEmail(ctx context.Context, in *peerpb.UserTokenReq) (*peerpb.UserTokenResp, error) {
	password, email := in.AccountMap[peerpb.AccountType_Password.String()], in.AccountMap[peerpb.AccountType_Email.String()]
	// 通过邮箱获取用户信息
	user := &usermodel.User{}
	err := l.svcCtx.UserCollection.Find(context.Background(), bson.M{"accountMap." + peerpb.AccountType_Email.String(): email}).One(user)
	if err != nil {
		l.Errorf("login by password email error: %v", err)
		return &peerpb.UserTokenResp{
			Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.PasswordInvalid),
		}, nil
	}
	// 验证密码
	ok := utils.Pwd.VerifyPwd(password, user.GetAccountMap().Get(peerpb.AccountType_Password.String()), user.GetAccountMap().Get(peerpb.AccountType_PasswordSalt.String()))
	if !ok {
		return &peerpb.UserTokenResp{
			Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.PasswordInvalid),
		}, nil
	}
	// 生成token
	return l.generateToken(in, user), nil
}

// LoginBySmsCode 手机号验证码登录
func (l *UserTokenLogic) LoginBySmsCode(ctx context.Context, in *peerpb.UserTokenReq) (*peerpb.UserTokenResp, error) {
	phone, phoneCountryCode := in.AccountMap[peerpb.AccountType_Phone.String()], in.AccountMap[peerpb.AccountType_PhoneCountryCode.String()]
	smsCode := in.VerifyMap[peerpb.AccountVerifyType_SmsCode.String()]
	// 验证验证码
	smsCodeVerifyResp, err := l.svcCtx.SmsService.SmsCodeVerify(context.Background(), &peerpb.SmsCodeVerifyReq{
		Header:    in.Header,
		Phone:     phone,
		PhoneCode: phoneCountryCode,
		Scene:     peerpb.SmsSceneType_SmsUserToken.String(),
		SmsCode:   smsCode,
		Delete:    false,
	})
	if err != nil {
		l.Errorf("login by sms code error: %v", err)
		return &peerpb.UserTokenResp{
			Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.SmsCodeInvalid),
		}, nil
	}
	if !smsCodeVerifyResp.Success {
		return &peerpb.UserTokenResp{
			Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.SmsCodeInvalid),
		}, nil
	}
	// 通过手机号获取用户信息
	user := &usermodel.User{}
	err = l.svcCtx.UserCollection.Find(context.Background(), bson.M{"accountMap." + peerpb.AccountType_Phone.String(): phone, "accountMap." + peerpb.AccountType_PhoneCountryCode.String(): phoneCountryCode}).One(user)
	if err != nil {
		l.Errorf("login by sms code error: %v", err)
		return &peerpb.UserTokenResp{
			Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.PhoneInvalid),
		}, nil
	}
	// 生成token
	return l.generateToken(in, user), nil
}

// LoginByEmailCode 邮箱登录
func (l *UserTokenLogic) LoginByEmailCode(ctx context.Context, in *peerpb.UserTokenReq) (*peerpb.UserTokenResp, error) {
	email := in.AccountMap[peerpb.AccountType_Email.String()]
	emailCode := in.VerifyMap[peerpb.AccountVerifyType_EmailCode.String()]
	// 验证验证码
	emailCodeVerifyResp, err := l.svcCtx.EmailService.EmailCodeVerify(context.Background(), &peerpb.EmailCodeVerifyReq{
		Header:    in.Header,
		Email:     email,
		Scene:     peerpb.EmailSceneType_EmailUserToken.String(),
		EmailCode: emailCode,
		Delete:    false,
	})
	if err != nil {
		l.Errorf("login by email code error: %v", err)
		return &peerpb.UserTokenResp{
			Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.EmailCodeInvalid),
		}, nil
	}
	if !emailCodeVerifyResp.Success {
		return &peerpb.UserTokenResp{
			Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.EmailCodeInvalid),
		}, nil
	}
	// 通过邮箱获取用户信息
	user := &usermodel.User{}
	err = l.svcCtx.UserCollection.Find(context.Background(), bson.M{"accountMap." + peerpb.AccountType_Email.String(): email}).One(user)
	if err != nil {
		l.Errorf("login by email code error: %v", err)
		return &peerpb.UserTokenResp{
			Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.EmailInvalid),
		}, nil
	}
	// 生成token
	return l.generateToken(in, user), nil
}

func (l *UserTokenLogic) generateToken(in *peerpb.UserTokenReq, user *usermodel.User) *peerpb.UserTokenResp {
	ssm := utils.SSM{
		"platform":    in.Header.Platform.String(),
		"installId":   in.Header.InstallId,
		"deviceModel": in.Header.DeviceModel,
	}
	var scope []string
	role := user.GetAccountMap().Get(peerpb.AccountType_Role.String())
	switch role {
	case usermodel.AccountRoleRobot:
		scope = []string{
			"^.*$", // TODO: 修改为机器人权限
		}
	default:
		scope = []string{
			"^.*$", // 表示所有权限
		}
	}
	tokenObject := l.svcCtx.Jwt.GenerateToken(user.UserId, usermodel.GetJwtUniqueKey(in.Header), int(user.GetAccountMap().GetInt64(peerpb.AccountType_Status.String())), ssm.Marshal(), scope)
	if role == usermodel.AccountRoleRobot {
		tokenObject.ExpiredAt = time.Now().AddDate(100, 0, 0).UnixMilli()
	}
	err := l.svcCtx.Jwt.SetToken(context.Background(), tokenObject)
	if err != nil {
		l.Errorf("set token error: %v", err)
		return &peerpb.UserTokenResp{
			Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.LoginFailed),
		}
	}
	l.Debugf("set token: %v", tokenObject)
	return &peerpb.UserTokenResp{
		UserId: user.UserId,
		Token:  tokenObject.Token,
	}
}
