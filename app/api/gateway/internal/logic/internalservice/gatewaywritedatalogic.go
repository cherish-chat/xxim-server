package internalservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/api/gateway/client/internalservice"
	"github.com/cherish-chat/xxim-server/app/api/gateway/internal/logic/connectionmanager"
	"google.golang.org/protobuf/proto"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/api/gateway/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GatewayWriteDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGatewayWriteDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GatewayWriteDataLogic {
	return &GatewayWriteDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GatewayWriteData 向用户推送数据
func (l *GatewayWriteDataLogic) GatewayWriteData(in *peerpb.GatewayWriteDataReq) (*peerpb.GatewayWriteDataResp, error) {
	if in.AllPods {
		return l.gatewayWriteDataToAllPods(in)
	}

	data, _ := proto.Marshal(in.Content)
	successConnections := make([]*peerpb.LongConnection, 0)
	if len(in.GetFilter().GetUserIds()) > 0 {
		connections := connectionmanager.ConnectionLogic.GetConnectionsByUserIds(in.GetFilter().GetUserIds())
		for _, connection := range connections {
			err := connection.SendMessage(context.Background(), data)
			if err != nil {
				l.Errorf("GatewayWriteData error: %v", err)
			} else {
				successConnections = append(successConnections, connection.ToPb())
			}
		}
	}
	return &peerpb.GatewayWriteDataResp{
		SuccessLongConnections: successConnections,
	}, nil
}

func (l *GatewayWriteDataLogic) gatewayWriteDataToAllPods(in *peerpb.GatewayWriteDataReq) (*peerpb.GatewayWriteDataResp, error) {
	var resp = &peerpb.GatewayWriteDataResp{}
	clients, err := l.svcCtx.Config.GetAllGatewayRpcClient()
	if err != nil {
		l.Errorf("GatewayWriteData error: %v", err)
		return resp, err
	}
	for _, client := range clients {
		resp, err := internalservice.NewInternalService(client).GatewayWriteData(context.Background(), &peerpb.GatewayWriteDataReq{
			Header:  in.Header,
			Filter:  in.Filter,
			Content: in.Content,
			AllPods: false,
		})
		if err != nil {
			l.Errorf("GatewayWriteData error: %v", err)
			return resp, err
		}
		resp.SuccessLongConnections = append(resp.SuccessLongConnections, resp.SuccessLongConnections...)
	}
	return resp, nil
}
