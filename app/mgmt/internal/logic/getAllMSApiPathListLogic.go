package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"
	"github.com/cherish-chat/xxim-server/common/xorm"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllMSApiPathListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllMSApiPathListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllMSApiPathListLogic {
	return &GetAllMSApiPathListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAllMSApiPathListLogic) GetAllMSApiPathList(in *pb.GetAllMSApiPathListReq) (*pb.GetAllMSApiPathListResp, error) {
	if in.Page == nil {
		in.Page = &pb.Page{Page: 1, Size: 0}
	}
	var models []*mgmtmodel.ApiPath
	wheres := xorm.NewGormWhere()
	if in.Filter != nil {
	}
	count, err := xorm.ListWithPagingOrder(l.svcCtx.Mysql(), &models, &mgmtmodel.ApiPath{}, in.Page.Page, in.Page.Size, "createTime DESC", wheres...)
	if err != nil {
		l.Errorf("GetAllMSApiPathList err: %v", err)
		return &pb.GetAllMSApiPathListResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var resp []*pb.MSApiPath
	for _, model := range models {
		role := model.ToPB()
		resp = append(resp, role)
	}
	return &pb.GetAllMSApiPathListResp{
		ApiPaths: resp,
		Total:    count,
	}, nil
}
