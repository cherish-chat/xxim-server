package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMSMenuDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMSMenuDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMSMenuDetailLogic {
	return &GetMSMenuDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMSMenuDetailLogic) GetMSMenuDetail(in *pb.GetMSMenuDetailReq) (*pb.GetMSMenuDetailResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetMSMenuDetailResp{}, nil
}
