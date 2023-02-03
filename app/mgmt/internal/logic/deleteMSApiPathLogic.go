package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMSApiPathLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMSApiPathLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMSApiPathLogic {
	return &DeleteMSApiPathLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteMSApiPathLogic) DeleteMSApiPath(in *pb.DeleteMSApiPathReq) (*pb.DeleteMSApiPathResp, error) {
	err := l.svcCtx.Mysql().Model(&mgmtmodel.ApiPath{}).Where("id in (?)", in.Ids).Delete(&mgmtmodel.ApiPath{}).Error
	if err != nil {
		l.Errorf("delete error: %v", err)
	}
	return &pb.DeleteMSApiPathResp{}, err
}
