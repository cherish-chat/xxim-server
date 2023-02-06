package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/xorm"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllUserInvitationCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllUserInvitationCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllUserInvitationCodeLogic {
	return &GetAllUserInvitationCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAllUserInvitationCodeLogic) GetAllUserInvitationCode(in *pb.GetAllUserInvitationCodeReq) (*pb.GetAllUserInvitationCodeResp, error) {
	if in.Page == nil {
		in.Page = &pb.Page{Page: 1, Size: 0}
	}
	var models []*usermodel.InvitationCode
	wheres := xorm.NewGormWhere()
	if in.Filter != nil {
	}
	count, err := xorm.ListWithPagingOrder(l.svcCtx.Mysql(), &models, &usermodel.InvitationCode{}, in.Page.Page, in.Page.Size, "createTime DESC", wheres...)
	if err != nil {
		l.Errorf("GetList err: %v", err)
		return &pb.GetAllUserInvitationCodeResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var resp []*pb.UserInvitationCode
	for _, model := range models {
		role := model.ToPB()
		resp = append(resp, role)
	}
	return &pb.GetAllUserInvitationCodeResp{
		UserInvitationCodes: resp,
		Total:               count,
	}, nil
}
