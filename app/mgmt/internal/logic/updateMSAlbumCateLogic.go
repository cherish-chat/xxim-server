package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"
	"time"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMSAlbumCateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateMSAlbumCateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMSAlbumCateLogic {
	return &UpdateMSAlbumCateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateMSAlbumCateLogic) UpdateMSAlbumCate(in *pb.UpdateMSAlbumCateReq) (*pb.UpdateMSAlbumCateResp, error) {
	// 查询原模型
	model := &mgmtmodel.AlbumCate{}
	err := l.svcCtx.Mysql().Model(model).Where("id = ?", in.AlbumCate.Id).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.UpdateMSAlbumCateResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	updateMap := make(map[string]interface{})
	updateMap["name"] = in.AlbumCate.Name
	updateMap["updateTime"] = time.Now().UnixMilli()
	if len(updateMap) > 0 {
		err = l.svcCtx.Mysql().Model(model).Where("id = ?", in.AlbumCate.Id).Updates(updateMap).Error
		if err != nil {
			l.Errorf("更新失败: %v", err)
			return &pb.UpdateMSAlbumCateResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	return &pb.UpdateMSAlbumCateResp{}, nil
}
