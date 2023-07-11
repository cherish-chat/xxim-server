package accountservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/i18n"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"go.mongodb.org/mongo-driver/bson"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserAccessTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserAccessTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAccessTokenLogic {
	return &UserAccessTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UserAccessToken 用户登录
func (l *UserAccessTokenLogic) UserAccessToken(in *pb.UserAccessTokenReq) (*pb.UserAccessTokenResp, error) {
	// todo: add your logic here and delete this line
	_, passwordExists := in.AccountMap[pb.AccountTypePassword]
	_, usernameExists := in.AccountMap[pb.AccountTypeUsername]
	_, phoneExists := in.AccountMap[pb.AccountTypePhone]
	_, phoneCountryCodeExists := in.AccountMap[pb.AccountTypePhoneCode]
	_, emailExists := in.AccountMap[pb.AccountTypeEmail]
	_, smsCodeExists := in.VerifyMap[pb.AccountVerifyTypeSmsCode]
	_, emailCodeExists := in.VerifyMap[pb.AccountVerifyTypeEmailCode]

	_, captchaIdExists := in.VerifyMap[pb.AccountVerifyTypeCaptchaId]
	_, captchaCodeExists := in.VerifyMap[pb.AccountVerifyTypeCaptchaCode]
	if !captchaIdExists || !captchaCodeExists {
		// 获取token是否需要图形验证码
		if l.svcCtx.Config.Account.Login.RequireCaptcha {
			// 参数错误
			return &pb.UserAccessTokenResp{
				Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, "参数错误, 登录时必须传入图形验证码"),
			}, nil
		}
	} else if captchaIdExists && captchaCodeExists {
		// 验证图形验证码
		captchaVerifyResp, err := l.svcCtx.CaptchaService.CaptchaVerify(l.ctx, &pb.CaptchaVerifyReq{
			Header:      in.Header,
			CaptchaId:   in.VerifyMap[pb.AccountVerifyTypeCaptchaId],
			CaptchaCode: in.VerifyMap[pb.AccountVerifyTypeCaptchaCode],
			Delete:      true,
		})
		if err != nil {
			return nil, err
		}
		if !captchaVerifyResp.Success {
			return &pb.UserAccessTokenResp{
				Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.CaptchaError),
			}, nil
		}
	}

	//平台
	if !utils.EnumInSlice[pb.Platform](in.Header.Platform, l.svcCtx.Config.Account.Register.AllowPlatform) {
		return &pb.UserAccessTokenResp{
			Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.PlatformNotAllow),
		}, nil
	}

	if passwordExists {
		// 密码登录，看看是用户名+密码、手机号+密码、还是邮箱+密码
		if usernameExists {
			return l.LoginByPasswordUsername(l.ctx, in)
		}
		if phoneExists && phoneCountryCodeExists {
			return l.LoginByPasswordPhone(l.ctx, in)
		}
		if emailExists {
			return l.LoginByPasswordEmail(l.ctx, in)
		}
		// 参数错误
		return &pb.UserAccessTokenResp{
			Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, "参数错误, 密码登录时必须传入用户名、手机号或邮箱"),
		}, nil
	} else {
		// 看看smsCode是否存在，存在则是手机号+验证码登录
		if smsCodeExists {
			// 判断手机号是否存在
			if phoneExists && phoneCountryCodeExists {
				return l.LoginBySmsCode(l.ctx, in)
			}
			// 参数错误
			return &pb.UserAccessTokenResp{
				Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, "参数错误, 手机号+验证码登录时必须传入手机号和国家代码"),
			}, nil
		}
		// 看看emailCode是否存在，存在则是邮箱+验证码登录
		if emailCodeExists {
			// 判断邮箱是否存在
			if emailExists {
				return l.LoginByEmailCode(l.ctx, in)
			}
			// 参数错误
			return &pb.UserAccessTokenResp{
				Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, "参数错误, 邮箱+验证码登录时必须传入邮箱"),
			}, nil
		}
		// 参数错误
		return &pb.UserAccessTokenResp{
			Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, "参数错误, 验证码登录时必须传入验证码"),
		}, nil
	}
}

// LoginByPasswordUsername 用户名密码登录
func (l *UserAccessTokenLogic) LoginByPasswordUsername(ctx context.Context, in *pb.UserAccessTokenReq) (*pb.UserAccessTokenResp, error) {
	password, username := in.AccountMap[pb.AccountTypePassword], in.AccountMap[pb.AccountTypeUsername]
	// 通过用户名获取用户信息
	user := &usermodel.User{}
	err := l.svcCtx.UserCollection.Find(l.ctx, bson.M{"accountMap." + pb.AccountTypeUsername: username}).One(user)
	if err != nil {
		l.Errorf("login by password username error: %v", err)
		return &pb.UserAccessTokenResp{
			Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.PasswordInvalid),
		}, nil
	}
	// 验证密码
	ok := utils.Pwd.VerifyPwd(password, user.GetAccountMap().Get(pb.AccountTypePassword), user.GetAccountMap().Get(pb.AccountTypePasswordSalt))
	if !ok {
		return &pb.UserAccessTokenResp{
			Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.PasswordInvalid),
		}, nil
	}
	// 生成token
	return l.generateToken(in, user), nil
}

// LoginByPasswordPhone 手机号密码登录
func (l *UserAccessTokenLogic) LoginByPasswordPhone(ctx context.Context, in *pb.UserAccessTokenReq) (*pb.UserAccessTokenResp, error) {
	password, phone, phoneCountryCode := in.AccountMap[pb.AccountTypePassword], in.AccountMap[pb.AccountTypePhone], in.AccountMap[pb.AccountTypePhoneCode]
	// 通过手机号获取用户信息
	user := &usermodel.User{}
	err := l.svcCtx.UserCollection.Find(l.ctx, bson.M{"accountMap." + pb.AccountTypePhone: phone, "accountMap." + pb.AccountTypePhoneCode: phoneCountryCode}).One(user)
	if err != nil {
		l.Errorf("login by password phone error: %v", err)
		return &pb.UserAccessTokenResp{
			Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.PasswordInvalid),
		}, nil
	}
	// 验证密码
	ok := utils.Pwd.VerifyPwd(password, user.GetAccountMap().Get(pb.AccountTypePassword), user.GetAccountMap().Get(pb.AccountTypePasswordSalt))
	if !ok {
		return &pb.UserAccessTokenResp{
			Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.PasswordInvalid),
		}, nil
	}
	// 生成token
	return l.generateToken(in, user), nil
}

// LoginByPasswordEmail 邮箱密码登录
func (l *UserAccessTokenLogic) LoginByPasswordEmail(ctx context.Context, in *pb.UserAccessTokenReq) (*pb.UserAccessTokenResp, error) {
	password, email := in.AccountMap[pb.AccountTypePassword], in.AccountMap[pb.AccountTypeEmail]
	// 通过邮箱获取用户信息
	user := &usermodel.User{}
	err := l.svcCtx.UserCollection.Find(l.ctx, bson.M{"accountMap." + pb.AccountTypeEmail: email}).One(user)
	if err != nil {
		l.Errorf("login by password email error: %v", err)
		return &pb.UserAccessTokenResp{
			Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.PasswordInvalid),
		}, nil
	}
	// 验证密码
	ok := utils.Pwd.VerifyPwd(password, user.GetAccountMap().Get(pb.AccountTypePassword), user.GetAccountMap().Get(pb.AccountTypePasswordSalt))
	if !ok {
		return &pb.UserAccessTokenResp{
			Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.PasswordInvalid),
		}, nil
	}
	// 生成token
	return l.generateToken(in, user), nil
}

// LoginBySmsCode 手机号验证码登录
func (l *UserAccessTokenLogic) LoginBySmsCode(ctx context.Context, in *pb.UserAccessTokenReq) (*pb.UserAccessTokenResp, error) {
	phone, phoneCountryCode := in.AccountMap[pb.AccountTypePhone], in.AccountMap[pb.AccountTypePhoneCode]
	smsCode := in.VerifyMap[pb.AccountVerifyTypeSmsCode]
	// 验证验证码
	smsCodeVerifyResp, err := l.svcCtx.SmsService.SmsCodeVerify(l.ctx, &pb.SmsCodeVerifyReq{
		Header:    in.Header,
		Phone:     phone,
		PhoneCode: phoneCountryCode,
		Scene:     pb.SmsSceneTypeUserToken,
		SmsCode:   smsCode,
		Delete:    false,
	})
	if err != nil {
		l.Errorf("login by sms code error: %v", err)
		return &pb.UserAccessTokenResp{
			Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.SmsCodeInvalid),
		}, nil
	}
	if !smsCodeVerifyResp.Success {
		return &pb.UserAccessTokenResp{
			Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.SmsCodeInvalid),
		}, nil
	}
	// 通过手机号获取用户信息
	user := &usermodel.User{}
	err = l.svcCtx.UserCollection.Find(l.ctx, bson.M{"accountMap." + pb.AccountTypePhone: phone, "accountMap." + pb.AccountTypePhoneCode: phoneCountryCode}).One(user)
	if err != nil {
		l.Errorf("login by sms code error: %v", err)
		return &pb.UserAccessTokenResp{
			Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.PhoneInvalid),
		}, nil
	}
	// 生成token
	return l.generateToken(in, user), nil
}

// LoginByEmailCode 邮箱登录
func (l *UserAccessTokenLogic) LoginByEmailCode(ctx context.Context, in *pb.UserAccessTokenReq) (*pb.UserAccessTokenResp, error) {
	email := in.AccountMap[pb.AccountTypeEmail]
	emailCode := in.VerifyMap[pb.AccountVerifyTypeEmailCode]
	// 验证验证码
	emailCodeVerifyResp, err := l.svcCtx.EmailService.EmailCodeVerify(l.ctx, &pb.EmailCodeVerifyReq{
		Header:    in.Header,
		Email:     email,
		Scene:     pb.EmailSceneTypeUserToken,
		EmailCode: emailCode,
		Delete:    false,
	})
	if err != nil {
		l.Errorf("login by email code error: %v", err)
		return &pb.UserAccessTokenResp{
			Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.EmailCodeInvalid),
		}, nil
	}
	if !emailCodeVerifyResp.Success {
		return &pb.UserAccessTokenResp{
			Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.EmailCodeInvalid),
		}, nil
	}
	// 通过邮箱获取用户信息
	user := &usermodel.User{}
	err = l.svcCtx.UserCollection.Find(l.ctx, bson.M{"accountMap." + pb.AccountTypeEmail: email}).One(user)
	if err != nil {
		l.Errorf("login by email code error: %v", err)
		return &pb.UserAccessTokenResp{
			Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.EmailInvalid),
		}, nil
	}
	// 生成token
	return l.generateToken(in, user), nil
}

func (l *UserAccessTokenLogic) generateToken(in *pb.UserAccessTokenReq, user *usermodel.User) *pb.UserAccessTokenResp {
	ssm := utils.SSM{
		"platform":    in.Header.Platform.ToString(),
		"installId":   in.Header.InstallId,
		"deviceModel": in.Header.DeviceModel,
	}
	var scope []string
	role := user.GetAccountMap().Get(pb.AccountTypeRole)
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
	tokenObject := l.svcCtx.Jwt.GenerateToken(user.UserId, in.Header.GetJwtUniqueKey(), int(user.GetAccountMap().GetInt64(pb.AccountTypeStatus)), ssm.Marshal(), scope)
	if role == usermodel.AccountRoleRobot {
		tokenObject.ExpiredAt = time.Now().AddDate(100, 0, 0).UnixMilli()
	}
	err := l.svcCtx.Jwt.SetToken(l.ctx, tokenObject)
	if err != nil {
		l.Errorf("set token error: %v", err)
		return &pb.UserAccessTokenResp{
			Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.LoginFailed),
		}
	}
	l.Debugf("set token: %v", tokenObject)
	return &pb.UserAccessTokenResp{
		UserId:      user.UserId,
		AccessToken: tokenObject.Token,
	}
}
