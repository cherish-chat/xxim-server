package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMSUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMSUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMSUserLogic {
	return &DeleteMSUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteMSUserLogic) DeleteMSUser(in *pb.DeleteMSUserReq) (*pb.DeleteMSUserResp, error) {
	// todo: add your logic here and delete this line

	return &pb.DeleteMSUserResp{}, nil
}
