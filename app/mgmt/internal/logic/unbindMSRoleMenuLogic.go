package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnbindMSRoleMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnbindMSRoleMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnbindMSRoleMenuLogic {
	return &UnbindMSRoleMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnbindMSRoleMenuLogic) UnbindMSRoleMenu(in *pb.UnbindMSRoleMenuReq) (*pb.UnbindMSRoleMenuResp, error) {
	// todo: add your logic here and delete this line

	return &pb.UnbindMSRoleMenuResp{}, nil
}
