package gatewayservicelogic

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

// GatewayBatchGetUserConnection 获取用户的连接
// 二次开发人员建议不修改此处逻辑
func (l *GatewayBatchGetUserConnectionLogic) GatewayBatchGetUserConnection(in *pb.GatewayBatchGetUserConnectionReq) (*pb.GatewayBatchGetUserConnectionResp, error) {
	connections := ConnectionLogic.GetConnectionsByUserIds(in.UserIds)
	var resp = &pb.GatewayBatchGetUserConnectionResp{}
	for _, connection := range connections {
		resp.Connections = append(resp.Connections, connection.ToPb())
	}
	return resp, nil
}
