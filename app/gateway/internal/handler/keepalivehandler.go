package handler

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/logic"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/common/i18n"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

func KeepAliveHandler(svcCtx *svc.ServiceContext) func(ctx context.Context, connection *logic.WsConnection, c *pb.GatewayApiRequest) (pb.ResponseCode, []byte, error) {
	return func(ctx context.Context, connection *logic.WsConnection, c *pb.GatewayApiRequest) (pb.ResponseCode, []byte, error) {
		var response *pb.GatewayApiResponse
		var err error
		response, err = logic.NewGatewayKeepAliveLogic(ctx, svcCtx).KeepAlive(connection, c)
		if err != nil {
			logx.WithContext(ctx).Errorf("logic.NewGatewayKeepAliveLogic(ctx, svcCtx).KeepAlive: %v", err)
			if response == nil {
				response = &pb.GatewayApiResponse{
					Header: i18n.NewServerError(connection.Header),
					Body:   nil,
				}
			}
			return response.Header.Code, MarshalResponse(connection.Header, response), nil
		}
		if response == nil {
			response = &pb.GatewayApiResponse{
				Header: i18n.NewOkHeader(),
				Body:   nil,
			}
		}
		if response.GetHeader() == nil {
			response.Header = i18n.NewOkHeader()
		}
		return response.Header.Code, MarshalResponse(connection.Header, response), nil
	}
}
