package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xpwd"
	"gorm.io/gorm"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserPasswordLogic {
	return &UpdateUserPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserPasswordLogic) UpdateUserPassword(in *pb.UpdateUserPasswordReq) (*pb.UpdateUserPasswordResp, error) {
	err := usermodel.FlushUserCache(l.ctx, l.svcCtx.Redis(), []string{in.CommonReq.UserId})
	if err != nil {
		l.Errorf("flush user cache failed, err: %v", err)
		return &pb.UpdateUserPasswordResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	user := &usermodel.User{}
	err = l.svcCtx.Mysql().Model(user).Where("id = ?", in.CommonReq.UserId).First(user).Error
	if err != nil {
		l.Errorf("get user failed, err: %v", err)
		return &pb.UpdateUserPasswordResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	// 对比旧密码
	if !xpwd.VerifyPwd(in.OldPassword, user.Password, user.PasswordSalt) {
		return &pb.UpdateUserPasswordResp{CommonResp: pb.NewAlertErrorResp("操作失败", "密码错误")}, nil
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
		return &pb.UpdateUserPasswordResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.UpdateUserPasswordResp{}, nil
}
