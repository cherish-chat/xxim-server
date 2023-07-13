package connectionservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/api/gateway/client/connectionservice"
	"github.com/cherish-chat/xxim-server/app/api/gateway/internal/logic/connectionmanager"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/api/gateway/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GatewayKickLongConnectionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGatewayKickLongConnectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GatewayKickLongConnectionLogic {
	return &GatewayKickLongConnectionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GatewayKickLongConnection 踢掉连接
func (l *GatewayKickLongConnectionLogic) GatewayKickLongConnection(in *peerpb.GatewayKickLongConnectionReq) (*peerpb.GatewayKickLongConnectionResp, error) {
	if in.AllPods {
		return l.gatewayKickLongConnectionAllPods(in)
	}

	if len(in.GetFilter().GetUserIds()) > 0 {
		connections := connectionmanager.ConnectionLogic.GetConnectionsByUserIds(in.GetFilter().GetUserIds())
		for _, connection := range connections {
			connection.Connection.CloseConnection()
		}
	}
	return &peerpb.GatewayKickLongConnectionResp{}, nil
}

func (l *GatewayKickLongConnectionLogic) gatewayKickLongConnectionAllPods(in *peerpb.GatewayKickLongConnectionReq) (*peerpb.GatewayKickLongConnectionResp, error) {
	clients, err := l.svcCtx.Config.GetAllGatewayRpcClient()
	if err != nil {
		l.Errorf("get all gateway rpc client error: %v", err)
		return &peerpb.GatewayKickLongConnectionResp{}, err
	}
	for _, client := range clients {
		_, err := connectionservice.NewConnectionService(client).GatewayKickLongConnection(context.Background(), &peerpb.GatewayKickLongConnectionReq{
			Header: in.Header,
			Filter: in.Filter,
		})
		if err != nil {
			l.Errorf("gateway kick long connection error: %v", err)
			return &peerpb.GatewayKickLongConnectionResp{}, err
		}
	}
	return &peerpb.GatewayKickLongConnectionResp{}, nil
}
