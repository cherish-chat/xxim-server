package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMyMSMenuListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMyMSMenuListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMyMSMenuListLogic {
	return &GetMyMSMenuListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMyMSMenuListLogic) GetMyMSMenuList(in *pb.GetMyMSMenuListReq) (*pb.GetMyMSMenuListResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetMyMSMenuListResp{}, nil
}
