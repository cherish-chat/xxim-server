package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"
	"time"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddMSApiPathLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddMSApiPathLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddMSApiPathLogic {
	return &AddMSApiPathLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddMSApiPathLogic) AddMSApiPath(in *pb.AddMSApiPathReq) (*pb.AddMSApiPathResp, error) {
	model := &mgmtmodel.ApiPath{
		Id:         mgmtmodel.GetId(l.svcCtx.Mysql(), &mgmtmodel.ApiPath{}, 10000),
		Title:      in.ApiPath.Title,
		Path:       in.ApiPath.Path,
		CreateTime: time.Now().UnixMilli(),
		UpdateTime: time.Now().UnixMilli(),
	}
	err := l.svcCtx.Mysql().Model(model).Create(model).Error
	if err != nil {
		l.Errorf("新增失败: %v", err)
		return &pb.AddMSApiPathResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.AddMSApiPathResp{}, nil
}
