package gatewayservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
	"google.golang.org/protobuf/proto"

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
	data, _ := proto.Marshal(in.Data)
	successConnections := make([]*pb.LongConnection, 0)
	if len(in.GetFilter().GetUserIds()) > 0 {
		connections := ConnectionLogic.GetConnectionsByUserIds(in.GetFilter().GetUserIds())
		for _, connection := range connections {
			err := connection.SendMessage(l.ctx, data)
			if err != nil {
				l.Errorf("GatewayWriteDataToWs error: %v", err)
			} else {
				successConnections = append(successConnections, connection.ToPb())
			}
		}
	}
	return &pb.GatewayWriteDataToWsResp{
		SuccessConnections: successConnections,
	}, nil
}
