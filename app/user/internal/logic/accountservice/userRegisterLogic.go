package accountservicelogic

import (
	"context"
	"fmt"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/i18n"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xcache"
	"github.com/cherish-chat/xxim-server/common/xmq"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"net/url"
	"time"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

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
func (l *UserRegisterLogic) UserRegister(in *pb.UserRegisterReq) (*pb.UserRegisterResp, error) {
	user := &usermodel.User{
		UserId:       in.UserId,
		RegisterTime: primitive.NewDateTimeFromTime(time.Now()),
		DestroyTime:  0,
		AccountMap: bson.M{
			pb.AccountTypeStatus: usermodel.AccountStatusNormal,
			pb.AccountTypeRole:   usermodel.AccountRoleUser,
		},
		Nickname:   "",
		Avatar:     "",
		ProfileMap: utils.Map.SS2SA(in.ProfileMap),
		ExtraMap:   utils.Map.SS2SA(in.ExtraMap),
	}
	//验证请求
	{
		//平台
		if !utils.EnumInSlice[pb.Platform](in.Header.Platform, l.svcCtx.Config.Account.Register.AllowPlatform) {
			return &pb.UserRegisterResp{
				Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.PlatformNotAllow),
			}, nil
		}

		//是否必填password
		username, ok := in.AccountMap[pb.AccountTypeUsername]
		if !ok {
			if l.svcCtx.Config.Account.Register.RequirePassword {
				return &pb.UserRegisterResp{
					Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.UsernameRequired),
				}, nil
			}
		} else {
			user.AccountMap[pb.AccountTypeUsername] = username
			if l.svcCtx.Config.Account.UserRegex != "" {
				if !utils.Regex.Match(l.svcCtx.Config.Account.UserRegex, username) {
					return &pb.UserRegisterResp{
						Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.UsernameFormatError),
					}, nil
				}
			}
		}
		if l.svcCtx.Config.Account.UsernameUnique {
			//用户名上锁
			ok, err := xcache.Lock.Lock(l.ctx, l.svcCtx.Redis, xcache.RedisVal.LockKeyUserUsername(username), 5)
			if err != nil || !ok {
				return &pb.UserRegisterResp{
					Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.UsernameLockError),
				}, nil
			}
			defer xcache.Lock.Unlock(l.ctx, l.svcCtx.Redis, xcache.RedisVal.LockKeyUserUsername(username))
			//检查用户名是否已存在
			found := &usermodel.User{}
			err = l.svcCtx.UserCollection.Find(l.ctx, bson.M{
				"accountMap." + pb.AccountTypeUsername: username,
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
				return &pb.UserRegisterResp{
					Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.UsernameAlreadyExists),
				}, nil
			}
		}
		passwordSalt, ok := in.AccountMap[pb.AccountTypePasswordSalt]
		if !ok {
			if l.svcCtx.Config.Account.Register.RequirePassword {
				return &pb.UserRegisterResp{
					Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.PasswordSaltRequired),
				}, nil
			}
		} else {
			user.AccountMap[pb.AccountTypePasswordSalt] = passwordSalt
		}
		password, ok := in.AccountMap[pb.AccountTypePassword]
		if !ok {
			if l.svcCtx.Config.Account.Register.RequirePassword {
				return &pb.UserRegisterResp{
					Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.PasswordRequired),
				}, nil
			}
		} else {
			user.AccountMap[pb.AccountTypePassword] = utils.Pwd.GeneratePwd(password, passwordSalt)
		}

		//是否必填手机号
		phone, ok := in.AccountMap[pb.AccountTypePhone]
		if !ok {
			if l.svcCtx.Config.Account.Register.RequireBindPhone {
				return &pb.UserRegisterResp{
					Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.PhoneRequired),
				}, nil
			}
		} else {
			user.AccountMap[pb.AccountTypePhone] = phone
		}
		phoneCode, ok := in.AccountMap[pb.AccountTypePhoneCode]
		if !ok {
			if l.svcCtx.Config.Account.Register.RequireBindPhone {
				return &pb.UserRegisterResp{
					Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.PhoneCodeRequired),
				}, nil
			}
		} else {
			user.AccountMap[pb.AccountTypePhoneCode] = phoneCode
		}
		if phone != "" && phoneCode != "" {
			if l.svcCtx.Config.Account.PhoneRegex != "" {
				if !utils.Regex.Match(l.svcCtx.Config.Account.PhoneRegex, phone) {
					return &pb.UserRegisterResp{
						Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.PhoneFormatError),
					}, nil
				}
			}
			if !utils.AnyInSlice[string](phoneCode, l.svcCtx.Config.Account.PhoneCode) {
				return &pb.UserRegisterResp{
					Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.PhoneCodeError),
				}, nil
			}
		}
		smsCode, ok := in.VerifyMap[pb.AccountVerifyTypeSmsCode]
		if !ok {
			if l.svcCtx.Config.Account.Register.RequireBindPhone {
				return &pb.UserRegisterResp{
					Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.SmsCodeRequired),
				}, nil
			}
		} else {
			//验证短信验证码
			smsVerifyResp, err := l.svcCtx.SmsService.SmsCodeVerify(l.ctx, &pb.SmsCodeVerifyReq{
				Header:    in.Header,
				Phone:     phone,
				PhoneCode: phoneCode,
				SmsCode:   smsCode,
				Delete:    true,
				Scene:     pb.SmsSceneTypeRegister,
			})
			if err != nil {
				l.Errorf("SmsVerify err: %v", err)
				return &pb.UserRegisterResp{
					Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.SmsCodeError),
				}, nil
			}
			if !smsVerifyResp.Success {
				return &pb.UserRegisterResp{
					Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.SmsCodeError),
				}, nil
			}
		}
		if l.svcCtx.Config.Account.PhoneUnique {
			//手机号上锁
			ok, err := xcache.Lock.Lock(l.ctx, l.svcCtx.Redis, xcache.RedisVal.LockKeyUserPhone(phone, phoneCode), 5)
			if err != nil || !ok {
				return &pb.UserRegisterResp{
					Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.PhoneLockError),
				}, nil
			}
			defer xcache.Lock.Unlock(l.ctx, l.svcCtx.Redis, xcache.RedisVal.LockKeyUserPhone(phone, phoneCode))
			//检查用户名是否已存在
			found := &usermodel.User{}
			err = l.svcCtx.UserCollection.Find(l.ctx, bson.M{
				"accountMap." + pb.AccountTypePhone:     phone,
				"accountMap." + pb.AccountTypePhoneCode: phoneCode,
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
				return &pb.UserRegisterResp{
					Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.PhoneAlreadyExists),
				}, nil
			}
		}

		//是否必填邮箱
		email, ok := in.AccountMap[pb.AccountTypeEmail]
		if !ok {
			if l.svcCtx.Config.Account.Register.RequireBindEmail {
				return &pb.UserRegisterResp{
					Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.EmailRequired),
				}, nil
			}
		} else {
			user.AccountMap[pb.AccountTypeEmail] = email
		}
		emailCode, ok := in.VerifyMap[pb.AccountVerifyTypeEmailCode]
		if !ok {
			if l.svcCtx.Config.Account.Register.RequireBindEmail {
				return &pb.UserRegisterResp{
					Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.EmailCodeInvalid),
				}, nil
			}
		} else {
			if l.svcCtx.Config.Account.EmailRegex != "" {
				if !utils.Regex.Match(l.svcCtx.Config.Account.EmailRegex, email) {
					return &pb.UserRegisterResp{
						Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.EmailFormatError),
					}, nil
				}
			}
			//验证邮箱验证码
			emailVerifyResp, err := l.svcCtx.EmailService.EmailCodeVerify(l.ctx, &pb.EmailCodeVerifyReq{
				Header:    in.Header,
				Email:     email,
				EmailCode: emailCode,
				Delete:    true,
				Scene:     pb.EmailSceneTypeRegister,
			})
			if err != nil {
				l.Errorf("EmailVerify err: %v", err)
				return &pb.UserRegisterResp{
					Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.EmailCodeError),
				}, nil
			}
			if !emailVerifyResp.Success {
				return &pb.UserRegisterResp{
					Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.EmailCodeError),
				}, nil
			}
		}
		if l.svcCtx.Config.Account.EmailUnique {
			//手机号上锁
			ok, err := xcache.Lock.Lock(l.ctx, l.svcCtx.Redis, xcache.RedisVal.LockKeyUserEmail(email), 5)
			if err != nil || !ok {
				return &pb.UserRegisterResp{
					Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.EmailLockError),
				}, nil
			}
			defer xcache.Lock.Unlock(l.ctx, l.svcCtx.Redis, xcache.RedisVal.LockKeyUserEmail(email))
			//检查用户名是否已存在
			found := &usermodel.User{}
			err = l.svcCtx.UserCollection.Find(l.ctx, bson.M{
				"accountMap." + pb.AccountTypeEmail: email,
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
				return &pb.UserRegisterResp{
					Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.EmailAlreadyExists),
				}, nil
			}
		}

		//是否必填昵称
		if in.Nickname == nil || *in.Nickname == "" {
			if l.svcCtx.Config.Account.Register.RequireNickname {
				return &pb.UserRegisterResp{
					Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.NicknameRequired),
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
				return &pb.UserRegisterResp{
					Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.AvatarRequired),
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
		captchaId, ok := in.VerifyMap[pb.AccountVerifyTypeCaptchaId]
		if !ok {
			if l.svcCtx.Config.Account.Register.RequireCaptcha {
				return &pb.UserRegisterResp{
					Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.CaptchaRequired),
				}, nil
			}
		}
		captchaCode, ok := in.VerifyMap[pb.AccountVerifyTypeCaptchaCode]
		if !ok {
			if l.svcCtx.Config.Account.Register.RequireCaptcha {
				return &pb.UserRegisterResp{
					Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.CaptchaRequired),
				}, nil
			}
		}
		if captchaId != "" && captchaCode != "" {
			//验证图形验证码
			captchaVerifyResp, err := l.svcCtx.CaptchaService.CaptchaVerify(l.ctx, &pb.CaptchaVerifyReq{
				Header:      in.Header,
				CaptchaId:   captchaId,
				CaptchaCode: captchaCode,
				Delete:      true,
			})
			if err != nil {
				l.Errorf("CaptchaVerify err: %v", err)
				return &pb.UserRegisterResp{
					Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.CaptchaError),
				}, nil
			}
			if !captchaVerifyResp.Success {
				return &pb.UserRegisterResp{
					Header: i18n.NewToastHeader(pb.ToastActionData_ERROR, i18n.CaptchaError),
				}, nil
			}
		}
	}

	_, err := l.svcCtx.UserCollection.InsertOne(l.ctx, user)
	if err != nil {
		l.Errorf("InsertOne err: %v", err)
		return nil, err
	}

	// afterRegister
	go l.svcCtx.MQ.Produce(context.Background(), xmq.TopicAfterRegister, []byte(user.UserId))
	return &pb.UserRegisterResp{}, nil
}
