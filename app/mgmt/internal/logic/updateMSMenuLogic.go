package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMSMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateMSMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMSMenuLogic {
	return &UpdateMSMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateMSMenuLogic) UpdateMSMenu(in *pb.UpdateMSMenuReq) (*pb.UpdateMSMenuResp, error) {
	// todo: add your logic here and delete this line

	return &pb.UpdateMSMenuResp{}, nil
}
