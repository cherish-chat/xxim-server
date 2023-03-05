package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMSAlbumCateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMSAlbumCateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMSAlbumCateLogic {
	return &DeleteMSAlbumCateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteMSAlbumCateLogic) DeleteMSAlbumCate(in *pb.DeleteMSAlbumCateReq) (*pb.DeleteMSAlbumCateResp, error) {
	err := l.svcCtx.Mysql().Model(&mgmtmodel.AlbumCate{}).Where("id in (?)", in.Ids).Delete(&mgmtmodel.AlbumCate{}).Error
	if err != nil {
		l.Errorf("delete error: %v", err)
	}
	return &pb.DeleteMSAlbumCateResp{}, err
}
