package gatewayservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GatewayGetConnectionByFilterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGatewayGetConnectionByFilterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GatewayGetConnectionByFilterLogic {
	return &GatewayGetConnectionByFilterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GatewayGetConnectionByFilter 通过条件获取用户的连接
func (l *GatewayGetConnectionByFilterLogic) GatewayGetConnectionByFilter(in *pb.GatewayGetConnectionByFilterReq) (*pb.GatewayGetConnectionByFilterResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GatewayGetConnectionByFilterResp{}, nil
}
