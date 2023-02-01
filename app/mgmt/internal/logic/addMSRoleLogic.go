package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddMSRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddMSRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddMSRoleLogic {
	return &AddMSRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddMSRoleLogic) AddMSRole(in *pb.AddMSRoleReq) (*pb.AddMSRoleResp, error) {
	// todo: add your logic here and delete this line

	return &pb.AddMSRoleResp{}, nil
}
