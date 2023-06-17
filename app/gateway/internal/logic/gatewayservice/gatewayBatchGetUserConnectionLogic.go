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

// GatewayBatchGetUserConnection 批量获取用户的连接
func (l *GatewayBatchGetUserConnectionLogic) GatewayBatchGetUserConnection(in *pb.GatewayBatchGetUserConnectionReq) (*pb.GatewayBatchGetUserConnectionResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GatewayBatchGetUserConnectionResp{}, nil
}
