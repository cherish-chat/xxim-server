package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMSRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMSRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMSRoleLogic {
	return &DeleteMSRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteMSRoleLogic) DeleteMSRole(in *pb.DeleteMSRoleReq) (*pb.DeleteMSRoleResp, error) {
	// todo: add your logic here and delete this line

	return &pb.DeleteMSRoleResp{}, nil
}
