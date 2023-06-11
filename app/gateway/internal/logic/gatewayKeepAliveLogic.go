package logic

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

func (l *GatewayKeepAliveLogic) GatewayKeepAlive(in *pb.GatewayKeepAliveReq) (*pb.GatewayKeepAliveResp, error) {
	return &pb.GatewayKeepAliveResp{}, status.Error(codes.Unimplemented, "cannot call method GatewayKeepAlive")
}

// KeepAlive 保持连接
// 客户端必须每隔 config.Websocket.KeepAliveSecond 秒发送一次心跳包
// 二次开发人员可以在这里修改逻辑，比如一致性算法安全校验等
func (l *GatewayKeepAliveLogic) KeepAlive(connection *WsConnection) (pb.ResponseCode, []byte, error) {
	WsManager.KeepAlive(connection)
	return pb.ResponseCode_SUCCESS, nil, nil
}
