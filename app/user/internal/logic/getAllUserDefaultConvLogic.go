package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/xorm"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllUserDefaultConvLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllUserDefaultConvLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllUserDefaultConvLogic {
	return &GetAllUserDefaultConvLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAllUserDefaultConvLogic) GetAllUserDefaultConv(in *pb.GetAllUserDefaultConvReq) (*pb.GetAllUserDefaultConvResp, error) {
	if in.Page == nil {
		in.Page = &pb.Page{Page: 1, Size: 0}
	}
	var models []*usermodel.DefaultConv
	wheres := xorm.NewGormWhere()
	if in.Filter != nil {
	}
	count, err := xorm.ListWithPagingOrder(l.svcCtx.Mysql(), &models, &usermodel.DefaultConv{}, in.Page.Page, in.Page.Size, "createTime DESC", wheres...)
	if err != nil {
		l.Errorf("GetList err: %v", err)
		return &pb.GetAllUserDefaultConvResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var resp []*pb.UserDefaultConv
	for _, model := range models {
		role := model.ToPB()
		resp = append(resp, role)
	}
	return &pb.GetAllUserDefaultConvResp{
		UserDefaultConvs: resp,
		Total:            count,
	}, nil
}
