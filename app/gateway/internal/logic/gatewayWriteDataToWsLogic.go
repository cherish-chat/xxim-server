package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GatewayWriteDataToWsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGatewayWriteDataToWsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GatewayWriteDataToWsLogic {
	return &GatewayWriteDataToWsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GatewayWriteDataToWs 向用户的连接写入数据
// 二次开发人员不建议修改此处逻辑
func (l *GatewayWriteDataToWsLogic) GatewayWriteDataToWs(in *pb.GatewayWriteDataToWsReq) (*pb.GatewayWriteDataToWsResp, error) {
	gatewayGetConnectionByFilterResp, err := NewGatewayGetConnectionByFilterLogic(l.ctx, l.svcCtx).GatewayGetConnectionByFilter(&pb.GatewayGetConnectionByFilterReq{
		Header: in.Header,
		Filter: in.Filter,
	})
	if err != nil {
		l.Errorf("GatewayGetConnectionByFilter error: %v", err)
		return &pb.GatewayWriteDataToWsResp{}, err
	}
	var resp = &pb.GatewayWriteDataToWsResp{}
	for _, connection := range gatewayGetConnectionByFilterResp.GetConnections() {
		success := WsManager.WriteData(connection.Id, in.Data)
		if success {
			resp.SuccessConnections = append(resp.SuccessConnections, connection)
		}
	}
	return resp, nil
}
