package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnbindMSUserRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnbindMSUserRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnbindMSUserRoleLogic {
	return &UnbindMSUserRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnbindMSUserRoleLogic) UnbindMSUserRole(in *pb.UnbindMSUserRoleReq) (*pb.UnbindMSUserRoleResp, error) {
	// todo: add your logic here and delete this line

	return &pb.UnbindMSUserRoleResp{}, nil
}
