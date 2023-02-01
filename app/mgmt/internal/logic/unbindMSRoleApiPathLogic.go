package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnbindMSRoleApiPathLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnbindMSRoleApiPathLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnbindMSRoleApiPathLogic {
	return &UnbindMSRoleApiPathLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnbindMSRoleApiPathLogic) UnbindMSRoleApiPath(in *pb.UnbindMSRoleApiPathReq) (*pb.UnbindMSRoleApiPathResp, error) {
	// todo: add your logic here and delete this line

	return &pb.UnbindMSRoleApiPathResp{}, nil
}
