package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type BindMSUserRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBindMSUserRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BindMSUserRoleLogic {
	return &BindMSUserRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BindMSUserRoleLogic) BindMSUserRole(in *pb.BindMSUserRoleReq) (*pb.BindMSUserRoleResp, error) {
	// todo: add your logic here and delete this line

	return &pb.BindMSUserRoleResp{}, nil
}
