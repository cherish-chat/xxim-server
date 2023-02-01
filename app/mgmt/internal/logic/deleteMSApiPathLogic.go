package logic

import (
	"context"

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
	// todo: add your logic here and delete this line

	return &pb.DeleteMSApiPathResp{}, nil
}
