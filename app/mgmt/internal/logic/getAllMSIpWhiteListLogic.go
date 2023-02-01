package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllMSIpWhiteListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllMSIpWhiteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllMSIpWhiteListLogic {
	return &GetAllMSIpWhiteListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAllMSIpWhiteListLogic) GetAllMSIpWhiteList(in *pb.GetAllMSIpWhiteListReq) (*pb.GetAllMSIpWhiteListResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetAllMSIpWhiteListResp{}, nil
}
