package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserInvitationCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserInvitationCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInvitationCodeLogic {
	return &UpdateUserInvitationCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateUserInvitationCodeLogic) UpdateUserInvitationCode(in *pb.UpdateUserInvitationCodeReq) (*pb.UpdateUserInvitationCodeResp, error) {
	// 查询原模型
	model := &usermodel.InvitationCode{}
	err := l.svcCtx.Mysql().Model(model).Where("code = ?", in.UserInvitationCode.Code).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.UpdateUserInvitationCodeResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	updateMap := map[string]interface{}{}
	{
		updateMap["remark"] = in.UserInvitationCode.Remark
		updateMap["isEnable"] = in.UserInvitationCode.IsEnable
		updateMap["defaultConvMode"] = in.UserInvitationCode.DefaultConvMode
	}
	if len(updateMap) > 0 {
		err = l.svcCtx.Mysql().Model(model).Where("code = ?", in.UserInvitationCode.Code).Updates(updateMap).Error
		if err != nil {
			l.Errorf("更新失败: %v", err)
			return &pb.UpdateUserInvitationCodeResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	return &pb.UpdateUserInvitationCodeResp{}, nil
}
