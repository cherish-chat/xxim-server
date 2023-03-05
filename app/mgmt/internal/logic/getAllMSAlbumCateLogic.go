package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllMSAlbumCateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllMSAlbumCateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllMSAlbumCateLogic {
	return &GetAllMSAlbumCateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 相册
func (l *GetAllMSAlbumCateLogic) GetAllMSAlbumCate(in *pb.GetAllMSAlbumCateReq) (*pb.GetAllMSAlbumCateResp, error) {
	var models []*mgmtmodel.AlbumCate
	err := l.svcCtx.Mysql().Model(&mgmtmodel.AlbumCate{}).Where("deleteTime = ?", 0).Find(&models).Error
	if err != nil {
		l.Errorf("select error: %v", err)
		return &pb.GetAllMSAlbumCateResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	var resp []*pb.MSAlbumCate
	for _, model := range models {
		resp = append(resp, model.ToPb())
	}
	return &pb.GetAllMSAlbumCateResp{
		AlbumCates: resp,
	}, nil
}
