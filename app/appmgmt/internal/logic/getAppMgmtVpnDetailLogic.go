package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAppMgmtVpnDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAppMgmtVpnDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppMgmtVpnDetailLogic {
	return &GetAppMgmtVpnDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAppMgmtVpnDetailLogic) GetAppMgmtVpnDetail(in *pb.GetAppMgmtVpnDetailReq) (*pb.GetAppMgmtVpnDetailResp, error) {
	// 查询原模型
	model := &appmgmtmodel.Vpn{}
	err := l.svcCtx.Mysql().Model(model).Where("id = ?", in.Id).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.GetAppMgmtVpnDetailResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.GetAppMgmtVpnDetailResp{AppMgmtVpn: model.ToPB()}, nil
}
