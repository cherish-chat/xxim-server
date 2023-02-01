package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllMSMenuListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllMSMenuListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllMSMenuListLogic {
	return &GetAllMSMenuListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAllMSMenuListLogic) GetAllMSMenuList(in *pb.GetAllMSMenuListReq) (*pb.GetAllMSMenuListResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetAllMSMenuListResp{}, nil
}
