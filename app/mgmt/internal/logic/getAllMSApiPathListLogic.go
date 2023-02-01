package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllMSApiPathListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllMSApiPathListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllMSApiPathListLogic {
	return &GetAllMSApiPathListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAllMSApiPathListLogic) GetAllMSApiPathList(in *pb.GetAllMSApiPathListReq) (*pb.GetAllMSApiPathListResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetAllMSApiPathListResp{}, nil
}
