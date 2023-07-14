package accountservicelogic

import (
	"context"
	"fmt"
	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xcache"
	"github.com/cherish-chat/xxim-server/common/xmq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/url"
	"time"

	"github.com/cherish-chat/xxim-server/app/service/user/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserRegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserRegisterLogic {
	return &UserRegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UserRegister 用户注册
func (l *UserRegisterLogic) UserRegister(in *peerpb.UserRegisterReq) (*peerpb.UserRegisterResp, error) {
	user := &usermodel.User{
		UserId:       in.UserId,
		RegisterTime: primitive.NewDateTimeFromTime(time.Now()),
		DestroyTime:  0,
		AccountMap: bson.M{
			peerpb.AccountType_Status.String(): usermodel.AccountStatusNormal,
			peerpb.AccountType_Role.String():   usermodel.AccountRoleUser,
		},
		Nickname:   "",
		Avatar:     "",
		ProfileMap: utils.Map.SS2SA(in.ProfileMap),
		ExtraMap:   utils.Map.SS2SA(in.ExtraMap),
	}
	//验证请求
	{
		//平台
		if !utils.EnumInSlice[peerpb.Platform](in.Header.Platform, l.svcCtx.Config.Account.Register.AllowPlatform) {
			return &peerpb.UserRegisterResp{
				Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.PlatformNotAllow),
			}, nil
		}

		//是否必填password
		username, ok := in.AccountMap[peerpb.AccountType_Username.String()]
		if !ok {
			if l.svcCtx.Config.Account.Register.RequireUsername {
				return &peerpb.UserRegisterResp{
					Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.UsernameRequired),
				}, nil
			}
		} else {
			user.AccountMap[peerpb.AccountType_Username.String()] = username
			if l.svcCtx.Config.Account.Register.UsernameRegex != "" {
				if !utils.Regex.Match(l.svcCtx.Config.Account.Register.UsernameRegex, username) {
					return &peerpb.UserRegisterResp{
						Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.UsernameFormatError),
					}, nil
				}
			}
		}
		if l.svcCtx.Config.Account.Register.UsernameUnique {
			//用户名上锁
			ok, err := xcache.Lock.Lock(context.Background(), l.svcCtx.Redis, xcache.RedisVal.LockKeyUserUsername(username), 5)
			if err != nil || !ok {
				return &peerpb.UserRegisterResp{
					Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.UsernameLockError),
				}, nil
			}
			defer xcache.Lock.Unlock(context.Background(), l.svcCtx.Redis, xcache.RedisVal.LockKeyUserUsername(username))
			//检查用户名是否已存在
			found := &usermodel.User{}
			err = l.svcCtx.UserCollection.Find(context.Background(), bson.M{
				"accountMap." + peerpb.AccountType_Username.String(): username,
			}).One(found)
			if err != nil {
				if err != mongo.ErrNoDocuments {
					l.Errorf("UserRegisterLogic.UserRegister l.svcCtx.UserCollection.Find error: %v", err)
					return nil, err
				} else {
					// 没问题
				}
			} else {
				// 已存在
				return &peerpb.UserRegisterResp{
					Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.UsernameAlreadyExists),
				}, nil
			}
		}
		passwordSalt, ok := in.AccountMap[peerpb.AccountType_PasswordSalt.String()]
		if !ok {
			if l.svcCtx.Config.Account.Register.RequirePassword {
				return &peerpb.UserRegisterResp{
					Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.PasswordSaltRequired),
				}, nil
			}
		} else {
			user.AccountMap[peerpb.AccountType_PasswordSalt.String()] = passwordSalt
		}
		password, ok := in.AccountMap[peerpb.AccountType_Password.String()]
		if !ok {
			if l.svcCtx.Config.Account.Register.RequirePassword {
				return &peerpb.UserRegisterResp{
					Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.PasswordRequired),
				}, nil
			}
		} else {
			user.AccountMap[peerpb.AccountType_Password.String()] = utils.Pwd.GeneratePwd(password, passwordSalt)
		}

		//是否必填手机号
		phone, ok := in.AccountMap[peerpb.AccountType_Phone.String()]
		if !ok {
			if l.svcCtx.Config.Account.Register.RequireBindPhone {
				return &peerpb.UserRegisterResp{
					Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.PhoneRequired),
				}, nil
			}
		} else {
			user.AccountMap[peerpb.AccountType_Phone.String()] = phone
		}
		phoneCode, ok := in.AccountMap[peerpb.AccountType_PhoneCountryCode.String()]
		if !ok {
			if l.svcCtx.Config.Account.Register.RequireBindPhone {
				return &peerpb.UserRegisterResp{
					Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.PhoneCodeRequired),
				}, nil
			}
		} else {
			user.AccountMap[peerpb.AccountType_PhoneCountryCode.String()] = phoneCode
		}
		if phone != "" && phoneCode != "" {
			foundRule := false
			for _, rule := range l.svcCtx.Config.Account.Register.PhoneRules {
				if rule.CountryCode == phoneCode {
					foundRule = true
					if rule.PhoneRegex != "" {
						if !utils.Regex.Match(rule.PhoneRegex, phone) {
							return &peerpb.UserRegisterResp{
								Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.PhoneFormatError),
							}, nil
						}
					}
					break
				}
			}
			if !foundRule {
				return &peerpb.UserRegisterResp{
					Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.PhoneCodeError),
				}, nil
			}
		}
		smsCode, ok := in.VerifyMap[peerpb.AccountVerifyType_SmsCode.String()]
		if !ok {
			if l.svcCtx.Config.Account.Register.RequireBindPhone {
				return &peerpb.UserRegisterResp{
					Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.SmsCodeRequired),
				}, nil
			}
		} else {
			//验证短信验证码
			smsVerifyResp, err := l.svcCtx.SmsService.SmsCodeVerify(context.Background(), &peerpb.SmsCodeVerifyReq{
				Header:    in.Header,
				Phone:     phone,
				PhoneCode: phoneCode,
				SmsCode:   smsCode,
				Delete:    true,
				Scene:     peerpb.SmsSceneType_SmsRegister.String(),
			})
			if err != nil {
				l.Errorf("SmsVerify err: %v", err)
				return &peerpb.UserRegisterResp{
					Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.SmsCodeError),
				}, nil
			}
			if !smsVerifyResp.Success {
				return &peerpb.UserRegisterResp{
					Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.SmsCodeError),
				}, nil
			}
		}
		if l.svcCtx.Config.Account.Register.PhoneUnique {
			//手机号上锁
			ok, err := xcache.Lock.Lock(context.Background(), l.svcCtx.Redis, xcache.RedisVal.LockKeyUserPhone(phone, phoneCode), 5)
			if err != nil || !ok {
				return &peerpb.UserRegisterResp{
					Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.PhoneLockError),
				}, nil
			}
			defer xcache.Lock.Unlock(context.Background(), l.svcCtx.Redis, xcache.RedisVal.LockKeyUserPhone(phone, phoneCode))
			//检查用户名是否已存在
			found := &usermodel.User{}
			err = l.svcCtx.UserCollection.Find(context.Background(), bson.M{
				"accountMap." + peerpb.AccountType_Phone.String():            phone,
				"accountMap." + peerpb.AccountType_PhoneCountryCode.String(): phoneCode,
			}).One(found)
			if err != nil {
				if err != mongo.ErrNoDocuments {
					l.Errorf("UserRegisterLogic.UserRegister l.svcCtx.UserCollection.Find error: %v", err)
					return nil, err
				} else {
					// 没问题
				}
			} else {
				// 已存在
				return &peerpb.UserRegisterResp{
					Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.PhoneAlreadyExists),
				}, nil
			}
		}

		//是否必填邮箱
		email, ok := in.AccountMap[peerpb.AccountType_Email.String()]
		if !ok {
			if l.svcCtx.Config.Account.Register.RequireBindEmail {
				return &peerpb.UserRegisterResp{
					Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.EmailRequired),
				}, nil
			}
		} else {
			user.AccountMap[peerpb.AccountType_Email.String()] = email
		}
		emailCode, ok := in.VerifyMap[peerpb.AccountVerifyType_EmailCode.String()]
		if !ok {
			if l.svcCtx.Config.Account.Register.RequireBindEmail {
				return &peerpb.UserRegisterResp{
					Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.EmailCodeInvalid),
				}, nil
			}
		} else {
			if l.svcCtx.Config.Account.Register.EmailRegex != "" {
				if !utils.Regex.Match(l.svcCtx.Config.Account.Register.EmailRegex, email) {
					return &peerpb.UserRegisterResp{
						Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.EmailFormatError),
					}, nil
				}
			}
			//验证邮箱验证码
			emailVerifyResp, err := l.svcCtx.EmailService.EmailCodeVerify(context.Background(), &peerpb.EmailCodeVerifyReq{
				Header:    in.Header,
				Email:     email,
				EmailCode: emailCode,
				Delete:    true,
				Scene:     peerpb.EmailSceneType_EmailRegister.String(),
			})
			if err != nil {
				l.Errorf("EmailVerify err: %v", err)
				return &peerpb.UserRegisterResp{
					Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.EmailCodeError),
				}, nil
			}
			if !emailVerifyResp.Success {
				return &peerpb.UserRegisterResp{
					Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.EmailCodeError),
				}, nil
			}
		}
		if l.svcCtx.Config.Account.Register.EmailUnique {
			//手机号上锁
			ok, err := xcache.Lock.Lock(context.Background(), l.svcCtx.Redis, xcache.RedisVal.LockKeyUserEmail(email), 5)
			if err != nil || !ok {
				return &peerpb.UserRegisterResp{
					Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.EmailLockError),
				}, nil
			}
			defer xcache.Lock.Unlock(context.Background(), l.svcCtx.Redis, xcache.RedisVal.LockKeyUserEmail(email))
			//检查用户名是否已存在
			found := &usermodel.User{}
			err = l.svcCtx.UserCollection.Find(context.Background(), bson.M{
				"accountMap." + peerpb.AccountType_Email.String(): email,
			}).One(found)
			if err != nil {
				if err != mongo.ErrNoDocuments {
					l.Errorf("UserRegisterLogic.UserRegister l.svcCtx.UserCollection.Find error: %v", err)
					return nil, err
				} else {
					// 没问题
				}
			} else {
				// 已存在
				return &peerpb.UserRegisterResp{
					Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.EmailAlreadyExists),
				}, nil
			}
		}

		//是否必填昵称
		if in.Nickname == nil || *in.Nickname == "" {
			if l.svcCtx.Config.Account.Register.RequireNickname {
				return &peerpb.UserRegisterResp{
					Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.NicknameRequired),
				}, nil
			} else {
				switch l.svcCtx.Config.Account.Register.DefaultNicknameRule {
				case "random":
					user.Nickname = fmt.Sprintf("%s%d", l.svcCtx.Config.Account.Register.RandomNicknamePrefix, utils.Random.Int(4))
				case "fixed":
					user.Nickname = l.svcCtx.Config.Account.Register.FixedNickname
				}
			}
		} else {
			user.Nickname = *in.Nickname
		}
		//是否必填头像
		if in.Avatar == nil || *in.Avatar == "" {
			if l.svcCtx.Config.Account.Register.RequireAvatar {
				return &peerpb.UserRegisterResp{
					Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.AvatarRequired),
				}, nil
			} else {
				switch l.svcCtx.Config.Account.Register.DefaultAvatarRule {
				case "byName":
					//根据昵称生成
					user.Avatar = fmt.Sprintf("/image/generateAvatar?text=%s&w=200&h=200&bg=%s&fg=%s", url.QueryEscape(user.Nickname), utils.Random.SliceString(l.svcCtx.Config.Account.Register.ByNameAvatarBgColors), utils.Random.SliceString(l.svcCtx.Config.Account.Register.ByNameAvatarFgColors))
				case "fixed":
					//固定
					user.Avatar = l.svcCtx.Config.Account.Register.FixedAvatar
				}
			}
		} else {
			user.Avatar = *in.Avatar
		}

		//是否验证图形验证码
		captchaId, ok := in.VerifyMap[peerpb.AccountVerifyType_CaptchaId.String()]
		if !ok {
			if l.svcCtx.Config.Account.Register.RequireCaptcha {
				return &peerpb.UserRegisterResp{
					Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.CaptchaRequired),
				}, nil
			}
		}
		captchaCode, ok := in.VerifyMap[peerpb.AccountVerifyType_CaptchaCode.String()]
		if !ok {
			if l.svcCtx.Config.Account.Register.RequireCaptcha {
				return &peerpb.UserRegisterResp{
					Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.CaptchaRequired),
				}, nil
			}
		}
		if captchaId != "" && captchaCode != "" {
			//验证图形验证码
			captchaVerifyResp, err := l.svcCtx.CaptchaService.CaptchaVerify(context.Background(), &peerpb.CaptchaVerifyReq{
				Header:      in.Header,
				CaptchaId:   captchaId,
				CaptchaCode: captchaCode,
				Delete:      true,
			})
			if err != nil {
				l.Errorf("CaptchaVerify err: %v", err)
				return &peerpb.UserRegisterResp{
					Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.CaptchaError),
				}, nil
			}
			if !captchaVerifyResp.Success {
				return &peerpb.UserRegisterResp{
					Header: peerpb.NewToastHeader(peerpb.ToastActionData_ERROR, peerpb.CaptchaError),
				}, nil
			}
		}
	}

	_, err := l.svcCtx.UserCollection.InsertOne(context.Background(), user)
	if err != nil {
		l.Errorf("InsertOne err: %v", err)
		return nil, err
	}

	// afterRegister
	go l.svcCtx.MQ.Produce(context.Background(), xmq.TopicAfterRegister, []byte(user.UserId))
	return &peerpb.UserRegisterResp{}, nil
}
