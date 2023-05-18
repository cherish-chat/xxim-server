package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GatewayGetUserConnectionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGatewayGetUserConnectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GatewayGetUserConnectionLogic {
	return &GatewayGetUserConnectionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GatewayGetUserConnectionLogic) GatewayGetUserConnection(in *pb.GatewayGetUserConnectionReq) (*pb.GatewayGetUserConnectionResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GatewayGetUserConnectionResp{}, nil
}
