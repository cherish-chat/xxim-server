package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindMSRoleApiPathLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBindMSRoleApiPathLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindMSRoleApiPathLogic {
	return &BindMSRoleApiPathLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BindMSRoleApiPathLogic) BindMSRoleApiPath(in *pb.BindMSRoleApiPathReq) (*pb.BindMSRoleApiPathResp, error) {
	// todo: add your logic here and delete this line

	return &pb.BindMSRoleApiPathResp{}, nil
}
