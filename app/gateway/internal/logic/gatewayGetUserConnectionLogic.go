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
	wsConnections, ok := WsManager.wsConnectionMap.GetByUserId(in.UserId)
	if !ok {
		return &pb.GatewayGetUserConnectionResp{}, nil
	}
	var resp = &pb.GatewayGetUserConnectionResp{}
	for _, wsConnection := range wsConnections {
		resp.Connections = append(resp.Connections, wsConnection.ToPb())
	}
	return resp, nil
}
