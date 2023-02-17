package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"time"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddMSAlbumCateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddMSAlbumCateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddMSAlbumCateLogic {
	return &AddMSAlbumCateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddMSAlbumCateLogic) AddMSAlbumCate(in *pb.AddMSAlbumCateReq) (*pb.AddMSAlbumCateResp, error) {
	model := &mgmtmodel.AlbumCate{
		Pid:        uint(utils.AnyToInt64(in.AlbumCate.Pid)),
		Type:       int(utils.AnyToInt64(in.AlbumCate.Type)),
		Name:       in.AlbumCate.Name,
		CreateTime: time.Now().UnixMilli(),
		UpdateTime: time.Now().UnixMilli(),
		DeleteTime: 0,
	}
	err := l.svcCtx.Mysql().Model(model).Create(model).Error
	if err != nil {
		l.Errorf("新增失败: %v", err)
		return &pb.AddMSAlbumCateResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.AddMSAlbumCateResp{}, nil
}
