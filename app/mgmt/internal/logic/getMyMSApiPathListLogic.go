package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMyMSApiPathListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMyMSApiPathListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMyMSApiPathListLogic {
	return &GetMyMSApiPathListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMyMSApiPathListLogic) GetMyMSApiPathList(in *pb.GetMyMSApiPathListReq) (*pb.GetMyMSApiPathListResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetMyMSApiPathListResp{}, nil
}
