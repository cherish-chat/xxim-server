package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xpwd"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"gorm.io/gorm"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ResetPasswordLogic) ResetPassword(in *pb.ResetPasswordReq) (*pb.ResetPasswordResp, error) {
	// 使用手机号查询用户是否存在
	var mobile = in.Mobile
	var mobileCountryCode = ""
	if in.MobileCountryCode != nil {
		mobileCountryCode = *in.MobileCountryCode
	}
	user := &usermodel.User{}
	err := l.svcCtx.Mysql().Model(user).Where("mobile = ? and mobileCountryCode = ?", mobile, mobileCountryCode).First(user).Error
	if err != nil {
		return &pb.ResetPasswordResp{}, nil
	}
	// 校验验证码
	var verifySmsResp *pb.VerifySmsResp
	xtrace.StartFuncSpan(l.ctx, "VerifySms", func(ctx context.Context) {
		verifySmsResp, err = NewVerifySmsLogic(l.ctx, l.svcCtx).VerifySms(&pb.VerifySmsReq{
			CommonReq:   in.CommonReq,
			Phone:       mobile,
			CountryCode: in.MobileCountryCode,
			Scene:       "resetPassword",
			Code:        in.SmsCode,
			Delete:      true,
		})
	})
	if err != nil {
		l.Errorf("verify sms failed, err: %v", err)
		return &pb.ResetPasswordResp{}, err
	}
	if verifySmsResp.GetCommonResp().Code != pb.CommonResp_Success {
		return &pb.ResetPasswordResp{CommonResp: verifySmsResp.CommonResp}, nil
	}
	err = usermodel.FlushUserCache(l.ctx, l.svcCtx.Redis(), []string{in.CommonReq.UserId})
	if err != nil {
		l.Errorf("flush user cache failed, err: %v", err)
		return &pb.ResetPasswordResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	// 更新用户信息
	updateMap := map[string]interface{}{}
	updateMap["password"] = xpwd.GeneratePwd(in.NewPassword, user.PasswordSalt)
	err = xorm.Transaction(l.svcCtx.Mysql(), func(tx *gorm.DB) error {
		err = tx.Model(user).Where("id = ?", in.CommonReq.UserId).Updates(updateMap).Error
		if err != nil {
			l.Errorf("update user failed, err: %v", err)
			return err
		}
		return nil
	}, func(tx *gorm.DB) error {
		err := usermodel.FlushUserCache(l.ctx, l.svcCtx.Redis(), []string{in.CommonReq.UserId})
		if err != nil {
			l.Errorf("flush user cache failed, err: %v", err)
			return err
		}
		return nil
	})
	if err != nil {
		return &pb.ResetPasswordResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.ResetPasswordResp{}, nil
}
