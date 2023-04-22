package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"
	"github.com/cherish-chat/xxim-server/common/xorm"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllAppMgmtLinkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllAppMgmtLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllAppMgmtLinkLogic {
	return &GetAllAppMgmtLinkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAllAppMgmtLinkLogic) GetAllAppMgmtLink(in *pb.GetAllAppMgmtLinkReq) (*pb.GetAllAppMgmtLinkResp, error) {
	if in.Page == nil {
		in.Page = &pb.Page{Page: 1, Size: 9999}
	}
	var models []*appmgmtmodel.Link
	wheres := xorm.NewGormWhere()
	if in.Filter != nil {
	}
	count, err := xorm.ListWithPagingOrder(l.svcCtx.Mysql(), &models, &appmgmtmodel.Link{}, in.Page.Page, in.Page.Size, "`sort` DESC", wheres...)
	if err != nil {
		l.Errorf("GetList err: %v", err)
		return &pb.GetAllAppMgmtLinkResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var resp []*pb.AppMgmtLink
	for _, model := range models {
		role := model.ToPB()
		resp = append(resp, role)
	}
	return &pb.GetAllAppMgmtLinkResp{
		AppMgmtLinks: resp,
		Total:        count,
	}, nil
}
