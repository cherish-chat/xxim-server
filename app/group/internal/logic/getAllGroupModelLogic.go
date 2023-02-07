package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/group/groupmodel"
	"github.com/cherish-chat/xxim-server/common/xorm"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllGroupModelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllGroupModelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllGroupModelLogic {
	return &GetAllGroupModelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetAllGroupModel 获取所有群组
func (l *GetAllGroupModelLogic) GetAllGroupModel(in *pb.GetAllGroupModelReq) (*pb.GetAllGroupModelResp, error) {
	if in.Page == nil {
		in.Page = &pb.Page{Page: 1, Size: 0}
	}
	var models []*groupmodel.Group
	wheres := xorm.NewGormWhere()
	if in.Filter != nil {
	}
	count, err := xorm.ListWithPagingOrder(l.svcCtx.Mysql(), &models, &groupmodel.Group{}, in.Page.Page, in.Page.Size, "createTime DESC", wheres...)
	if err != nil {
		l.Errorf("GetList err: %v", err)
		return &pb.GetAllGroupModelResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var resp []*pb.GroupModel
	for _, model := range models {
		role := model.ToPB()
		resp = append(resp, role)
	}
	return &pb.GetAllGroupModelResp{
		GroupModels: resp,
		Total:       count,
	}, nil
}
