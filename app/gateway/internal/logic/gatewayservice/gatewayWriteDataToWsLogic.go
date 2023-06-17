package gatewayservicelogic

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
func (l *GatewayWriteDataToWsLogic) GatewayWriteDataToWs(in *pb.GatewayWriteDataToWsReq) (*pb.GatewayWriteDataToWsResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GatewayWriteDataToWsResp{}, nil
}
