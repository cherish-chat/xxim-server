package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindMSRoleMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBindMSRoleMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindMSRoleMenuLogic {
	return &BindMSRoleMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BindMSRoleMenuLogic) BindMSRoleMenu(in *pb.BindMSRoleMenuReq) (*pb.BindMSRoleMenuResp, error) {
	// todo: add your logic here and delete this line

	return &pb.BindMSRoleMenuResp{}, nil
}
