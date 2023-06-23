package gatewayservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/gateway/client/gatewayservice"

	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GatewayWriteDataToWsWrapperLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGatewayWriteDataToWsWrapperLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GatewayWriteDataToWsWrapperLogic {
	return &GatewayWriteDataToWsWrapperLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GatewayWriteDataToWsWrapper 向用户的连接写入数据
func (l *GatewayWriteDataToWsWrapperLogic) GatewayWriteDataToWsWrapper(in *pb.GatewayWriteDataToWsWrapperReq) (*pb.GatewayWriteDataToWsResp, error) {
	pods, err := l.getGatewayPods()
	if err != nil {
		l.Errorf("get gateway pods error: %v", err)
		return &pb.GatewayWriteDataToWsResp{}, err
	}
	req := &pb.GatewayWriteDataToWsReq{
		Header: in.Header,
		Filter: in.Filter,
		Data:   in.Data,
	}
	for _, pod := range pods {
		gatewayWriteDataToWsResp, err := pod.GatewayWriteDataToWs(l.ctx, req)
		if err != nil {
			l.Errorf("gateway write data to ws error: %v", err)
			return &pb.GatewayWriteDataToWsResp{}, err
		}
		_ = gatewayWriteDataToWsResp
	}
	return &pb.GatewayWriteDataToWsResp{}, nil
}

func (l *GatewayWriteDataToWsWrapperLogic) getGatewayPods() ([]gatewayservice.GatewayService, error) {
	// todo 查询所有的gateway pod
	return []gatewayservice.GatewayService{
		l.svcCtx.GatewayService(),
	}, nil
}
