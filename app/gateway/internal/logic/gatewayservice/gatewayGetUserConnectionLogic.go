package gatewayservicelogic

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

// GatewayGetUserConnection 获取用户的连接
// 二次开发人员不建议修改此处逻辑
func (l *GatewayGetUserConnectionLogic) GatewayGetUserConnection(in *pb.GatewayGetUserConnectionReq) (*pb.GatewayGetUserConnectionResp, error) {
	connections := ConnectionLogic.GetConnectionsByUserIds([]string{in.UserId})
	if len(connections) == 0 {
		return &pb.GatewayGetUserConnectionResp{}, nil
	}
	var resp = &pb.GatewayGetUserConnectionResp{}
	for _, connection := range connections {
		resp.Connections = append(resp.Connections, connection.ToPb())
	}
	return resp, nil
}
