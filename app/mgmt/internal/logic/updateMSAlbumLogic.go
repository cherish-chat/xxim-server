package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"
	"time"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMSAlbumLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateMSAlbumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMSAlbumLogic {
	return &UpdateMSAlbumLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateMSAlbumLogic) UpdateMSAlbum(in *pb.UpdateMSAlbumReq) (*pb.UpdateMSAlbumResp, error) {
	// 查询原模型
	model := &mgmtmodel.Album{}
	err := l.svcCtx.Mysql().Model(model).Where("id = ?", in.Album.Id).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.UpdateMSAlbumResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	updateMap := make(map[string]interface{})
	updateMap["name"] = in.Album.Name
	updateMap["updateTime"] = time.Now().UnixMilli()
	if len(updateMap) > 0 {
		err = l.svcCtx.Mysql().Model(model).Where("id = ?", in.Album.Id).Updates(updateMap).Error
		if err != nil {
			l.Errorf("更新失败: %v", err)
			return &pb.UpdateMSAlbumResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	return &pb.UpdateMSAlbumResp{}, nil
}
