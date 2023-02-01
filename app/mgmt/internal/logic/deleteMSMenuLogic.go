package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMSMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMSMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMSMenuLogic {
	return &DeleteMSMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteMSMenuLogic) DeleteMSMenu(in *pb.DeleteMSMenuReq) (*pb.DeleteMSMenuResp, error) {
	// todo: add your logic here and delete this line

	return &pb.DeleteMSMenuResp{}, nil
}
