package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMSRoleDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMSRoleDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMSRoleDetailLogic {
	return &GetMSRoleDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMSRoleDetailLogic) GetMSRoleDetail(in *pb.GetMSRoleDetailReq) (*pb.GetMSRoleDetailResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetMSRoleDetailResp{}, nil
}
