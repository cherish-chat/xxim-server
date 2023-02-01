package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllMSRoleListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllMSRoleListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllMSRoleListLogic {
	return &GetAllMSRoleListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAllMSRoleListLogic) GetAllMSRoleList(in *pb.GetAllMSRoleListReq) (*pb.GetAllMSRoleListResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetAllMSRoleListResp{}, nil
}
