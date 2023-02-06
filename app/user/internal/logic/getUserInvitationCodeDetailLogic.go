package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInvitationCodeDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInvitationCodeDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInvitationCodeDetailLogic {
	return &GetUserInvitationCodeDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInvitationCodeDetailLogic) GetUserInvitationCodeDetail(in *pb.GetUserInvitationCodeDetailReq) (*pb.GetUserInvitationCodeDetailResp, error) {
	// 查询原模型
	model := &usermodel.InvitationCode{}
	err := l.svcCtx.Mysql().Model(model).Where("code = ?", in.Code).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.GetUserInvitationCodeDetailResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.GetUserInvitationCodeDetailResp{UserInvitationCode: model.ToPB()}, nil
}
