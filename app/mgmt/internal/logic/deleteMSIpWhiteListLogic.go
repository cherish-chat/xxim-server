package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMSIpWhiteListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMSIpWhiteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMSIpWhiteListLogic {
	return &DeleteMSIpWhiteListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteMSIpWhiteListLogic) DeleteMSIpWhiteList(in *pb.DeleteMSIpWhiteListReq) (*pb.DeleteMSIpWhiteListResp, error) {
	// todo: add your logic here and delete this line

	return &pb.DeleteMSIpWhiteListResp{}, nil
}
