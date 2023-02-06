package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteUserInvitationCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteUserInvitationCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteUserInvitationCodeLogic {
	return &DeleteUserInvitationCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteUserInvitationCodeLogic) DeleteUserInvitationCode(in *pb.DeleteUserInvitationCodeReq) (*pb.DeleteUserInvitationCodeResp, error) {
	model := &usermodel.InvitationCode{}
	err := l.svcCtx.Mysql().Model(model).Where("code in (?)", in.Codes).Delete(model).Error
	if err != nil {
		l.Errorf("delete error: %v", err)
		return &pb.DeleteUserInvitationCodeResp{
			CommonResp: pb.NewRetryErrorResp(),
		}, err
	}
	return &pb.DeleteUserInvitationCodeResp{}, nil
}
