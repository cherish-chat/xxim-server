package gatewayservicelogic

import (
	"context"
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

// GatewayKickWs 踢出用户的连接
// 二次开发人员可以在此处修改踢出用户连接的逻辑
// 比如踢出连接之前，先给用户发送一条消息
func (l *GatewayKickWsLogic) GatewayKickWs(in *pb.GatewayKickWsReq) (*pb.GatewayKickWsResp, error) {
	if len(in.GetFilter().GetUserIds()) > 0 {
		connections := ConnectionLogic.GetConnectionsByUserIds(in.GetFilter().GetUserIds())
		for _, connection := range connections {
			connection.Connection.CloseConnection(in.CloseCode, in.CloseReason)
		}
	}
	return &pb.GatewayKickWsResp{}, nil
}
