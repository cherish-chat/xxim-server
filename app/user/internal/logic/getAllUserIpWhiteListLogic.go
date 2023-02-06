package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/xorm"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllUserIpWhiteListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllUserIpWhiteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllUserIpWhiteListLogic {
	return &GetAllUserIpWhiteListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAllUserIpWhiteListLogic) GetAllUserIpWhiteList(in *pb.GetAllUserIpWhiteListReq) (*pb.GetAllUserIpWhiteListResp, error) {
	if in.Page == nil {
		in.Page = &pb.Page{Page: 1, Size: 0}
	}
	var models []*usermodel.IpWhiteList
	wheres := xorm.NewGormWhere()
	if in.Filter != nil {
	}
	count, err := xorm.ListWithPagingOrder(l.svcCtx.Mysql(), &models, &usermodel.IpWhiteList{}, in.Page.Page, in.Page.Size, "createTime DESC", wheres...)
	if err != nil {
		l.Errorf("GetList err: %v", err)
		return &pb.GetAllUserIpWhiteListResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var resp []*pb.UserIpList
	for _, model := range models {
		role := model.ToPB()
		resp = append(resp, role)
	}
	return &pb.GetAllUserIpWhiteListResp{
		UserIpLists: resp,
		Total:       count,
	}, nil
}
