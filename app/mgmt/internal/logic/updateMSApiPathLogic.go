package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMSApiPathLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateMSApiPathLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMSApiPathLogic {
	return &UpdateMSApiPathLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateMSApiPathLogic) UpdateMSApiPath(in *pb.UpdateMSApiPathReq) (*pb.UpdateMSApiPathResp, error) {
	// todo: add your logic here and delete this line

	return &pb.UpdateMSApiPathResp{}, nil
}
