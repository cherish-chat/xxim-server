package logic

import (
	"context"
	"nhooyr.io/websocket"

	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GatewayKickWsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGatewayKickWsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GatewayKickWsLogic {
	return &GatewayKickWsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GatewayKickWsLogic) GatewayKickWs(in *pb.GatewayKickWsReq) (*pb.GatewayKickWsResp, error) {
	gatewayGetConnectionByFilterResp, err := NewGatewayGetConnectionByFilterLogic(l.ctx, l.svcCtx).GatewayGetConnectionByFilter(&pb.GatewayGetConnectionByFilterReq{
		Header: in.Header,
		Filter: in.Filter,
	})
	if err != nil {
		l.Errorf("GatewayGetConnectionByFilter error: %v", err)
		return &pb.GatewayKickWsResp{}, err
	}
	for _, connection := range gatewayGetConnectionByFilterResp.GetConnections() {
		WsManager.CloseConnection(connection.Id, websocket.StatusCode(in.CloseCode.Number()), in.CloseReason)
	}
	return &pb.GatewayKickWsResp{}, nil
}
