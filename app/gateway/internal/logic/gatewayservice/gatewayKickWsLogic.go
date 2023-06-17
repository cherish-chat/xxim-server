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
func (l *GatewayKickWsLogic) GatewayKickWs(in *pb.GatewayKickWsReq) (*pb.GatewayKickWsResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GatewayKickWsResp{}, nil
}
