package connectionservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/api/gateway/client/connectionservice"
	"github.com/cherish-chat/xxim-server/app/api/gateway/internal/logic/connectionmanager"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/api/gateway/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLongConnectionLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLongConnectionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLongConnectionLogic {
	return &ListLongConnectionLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListLongConnection 获取长连接列表
func (l *ListLongConnectionLogic) ListLongConnection(in *peerpb.ListLongConnectionReq) (*peerpb.ListLongConnectionResp, error) {
	if in.AllPods {
		return l.listLongConnectionAllPods(in)
	}
	connections := connectionmanager.ConnectionLogic.GetConnectionsByUserIds(in.GetFilter().GetUserIds())
	var resp = &peerpb.ListLongConnectionResp{}
	if len(connections) == 0 {
		return resp, nil
	}
	for _, connection := range connections {
		resp.LongConnections = append(resp.LongConnections, connection.ToPb())
	}

	return resp, nil
}

func (l *ListLongConnectionLogic) listLongConnectionAllPods(in *peerpb.ListLongConnectionReq) (*peerpb.ListLongConnectionResp, error) {
	var resp = &peerpb.ListLongConnectionResp{}
	clients, err := l.svcCtx.Config.GetAllGatewayRpcClient()
	if err != nil {
		l.Errorf("get all gateway rpc client error: %v", err)
		return resp, err
	}
	for _, client := range clients {
		resp, err := connectionservice.NewConnectionService(client).ListLongConnection(context.Background(), &peerpb.ListLongConnectionReq{
			Header: in.Header,
			Filter: in.Filter,
		})
		if err != nil {
			l.Errorf("list long connection error: %v", err)
			return resp, err
		}
		resp.LongConnections = append(resp.LongConnections, resp.LongConnections...)
	}
	return resp, nil
}
