package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLatestVersionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLatestVersionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLatestVersionLogic {
	return &GetLatestVersionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetLatestVersionLogic) GetLatestVersion(in *pb.GetLatestVersionReq) (*pb.GetLatestVersionResp, error) {
	platform := in.CommonReq.Platform
	dest := &appmgmtmodel.Version{}
	err := l.svcCtx.Mysql().Model(dest).Where("platform = ?", platform).Order("createTime desc").First(dest).Error
	if err != nil {
		return &pb.GetLatestVersionResp{
			AppMgmtVersion: &pb.AppMgmtVersion{
				Version:  in.CommonReq.AppVersion,
				Platform: platform,
			},
		}, err
	}
	return &pb.GetLatestVersionResp{AppMgmtVersion: dest.ToPB()}, nil
}
