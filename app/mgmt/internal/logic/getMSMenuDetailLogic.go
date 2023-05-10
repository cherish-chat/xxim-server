package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMSMenuDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMSMenuDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMSMenuDetailLogic {
	return &GetMSMenuDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMSMenuDetailLogic) GetMSMenuDetail(in *pb.GetMSMenuDetailReq) (*pb.GetMSMenuDetailResp, error) {
	var model = &mgmtmodel.Menu{}
	err := l.svcCtx.Mysql().Model(&mgmtmodel.Menu{}).
		Where("id = ?", in.Id).
		First(model).Error
	if err != nil {
		l.Errorf("查询菜单失败: %v", err)
		return &pb.GetMSMenuDetailResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.GetMSMenuDetailResp{Menu: model.ToPb(func(key string) string {
		return l.svcCtx.T(in.CommonReq.Language, key)
	})}, nil
}
