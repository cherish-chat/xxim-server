package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/xpwd"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMSUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateMSUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMSUserLogic {
	return &UpdateMSUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateMSUserLogic) UpdateMSUser(in *pb.UpdateMSUserReq) (*pb.UpdateMSUserResp, error) {
	// 查询原模型
	model := &mgmtmodel.User{}
	err := l.svcCtx.Mysql().Model(model).Where("id = ?", in.User.Id).First(model).Error
	if err != nil {
		l.Errorf("查询用户失败: %v", err)
		return &pb.UpdateMSUserResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	updateMap := make(map[string]interface{})
	if in.User.Nickname != "" {
		updateMap["nickname"] = in.User.Nickname
	}
	if in.User.Avatar != "" {
		updateMap["avatar"] = in.User.Avatar
	}
	if in.User.Role != "" {
		updateMap["roleId"] = in.User.Role
	}
	if in.User.IsDisable {
		if !model.IsDisable {
			// 封禁
			updateMap["isDisable"] = true
		}
	} else {
		if model.IsDisable {
			// 解封
			updateMap["isDisable"] = false
		}
	}
	if in.User.Password != "" {
		updateMap["password"] = xpwd.GeneratePwd(in.User.Password, model.PasswordSalt)
	}
	if len(updateMap) > 0 {
		err = l.svcCtx.Mysql().Model(model).Where("id = ?", in.User.Id).Updates(updateMap).Error
		if err != nil {
			l.Errorf("更新用户失败: %v", err)
			return &pb.UpdateMSUserResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	return &pb.UpdateMSUserResp{}, nil
}
