package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMSAlbumLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMSAlbumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMSAlbumLogic {
	return &DeleteMSAlbumLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteMSAlbumLogic) DeleteMSAlbum(in *pb.DeleteMSAlbumReq) (*pb.DeleteMSAlbumResp, error) {
	err := l.svcCtx.Mysql().Model(&mgmtmodel.Album{}).Where("id in (?)", in.Ids).Delete(&mgmtmodel.Album{}).Error
	if err != nil {
		l.Errorf("delete error: %v", err)
	}
	return &pb.DeleteMSAlbumResp{}, err
}
