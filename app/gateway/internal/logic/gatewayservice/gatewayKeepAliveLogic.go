package gatewayservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GatewayKeepAliveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGatewayKeepAliveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GatewayKeepAliveLogic {
	return &GatewayKeepAliveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// KeepAlive 保持连接
func (l *GatewayKeepAliveLogic) GatewayKeepAlive(in *pb.GatewayKeepAliveReq) (*pb.GatewayKeepAliveResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GatewayKeepAliveResp{}, nil
}
