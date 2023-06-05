package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GatewayBatchGetUserConnectionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGatewayBatchGetUserConnectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GatewayBatchGetUserConnectionLogic {
	return &GatewayBatchGetUserConnectionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GatewayBatchGetUserConnectionLogic) GatewayBatchGetUserConnection(in *pb.GatewayBatchGetUserConnectionReq) (*pb.GatewayBatchGetUserConnectionResp, error) {
	wsConnections := WsManager.wsConnectionMap.GetByUserIds(in.UserIds)
	var resp = &pb.GatewayBatchGetUserConnectionResp{}
	for _, wsConnection := range wsConnections {
		resp.Connections = append(resp.Connections, wsConnection.ToPb())
	}
	return resp, nil
}
