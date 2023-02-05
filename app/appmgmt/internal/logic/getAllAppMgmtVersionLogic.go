package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"
	"github.com/cherish-chat/xxim-server/common/xorm"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllAppMgmtVersionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllAppMgmtVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllAppMgmtVersionLogic {
	return &GetAllAppMgmtVersionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAllAppMgmtVersionLogic) GetAllAppMgmtVersion(in *pb.GetAllAppMgmtVersionReq) (*pb.GetAllAppMgmtVersionResp, error) {
	if in.Page == nil {
		in.Page = &pb.Page{Page: 1, Size: 0}
	}
	var models []*appmgmtmodel.Version
	wheres := xorm.NewGormWhere()
	if in.Filter != nil {
	}
	count, err := xorm.ListWithPagingOrder(l.svcCtx.Mysql(), &models, &appmgmtmodel.Version{}, in.Page.Page, in.Page.Size, "createTime DESC", wheres...)
	if err != nil {
		l.Errorf("GetList err: %v", err)
		return &pb.GetAllAppMgmtVersionResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var resp []*pb.AppMgmtVersion
	for _, model := range models {
		role := model.ToPB()
		resp = append(resp, role)
	}
	return &pb.GetAllAppMgmtVersionResp{
		AppMgmtVersions: resp,
		Total:           count,
	}, nil
}
