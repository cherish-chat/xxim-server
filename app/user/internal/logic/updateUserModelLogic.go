package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/xpwd"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserModelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserModelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserModelLogic {
	return &UpdateUserModelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserModelLogic) UpdateUserModel(in *pb.UpdateUserModelReq) (*pb.UpdateUserModelResp, error) {
	// 查询原模型
	model := &usermodel.User{}
	err := l.svcCtx.Mysql().Model(model).Where("id = ?", in.UserModel.Id).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.UpdateUserModelResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	updateMap := map[string]interface{}{}
	{
		updateMap["mobile"] = in.UserModel.Mobile
		updateMap["mobileCountryCode"] = in.UserModel.MobileCountryCode
		updateMap["nickname"] = in.UserModel.Nickname
		updateMap["avatar"] = in.UserModel.Avatar
		updateMap["xb"] = in.UserModel.Xb
		updateMap["role"] = in.UserModel.Role
		updateMap["adminRemark"] = in.UserModel.AdminRemark
	}
	if in.Password != "" {
		updateMap["password"] = xpwd.GeneratePwd(in.Password, model.PasswordSalt)
	}
	if len(updateMap) > 0 {
		err = l.svcCtx.Mysql().Model(model).Where("id = ?", in.UserModel.Id).Updates(updateMap).Error
		if err != nil {
			l.Errorf("更新失败: %v", err)
			return &pb.UpdateUserModelResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	return &pb.UpdateUserModelResp{}, nil
}
