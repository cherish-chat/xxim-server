package connectionhandler

import (
	"context"
	gatewayservicelogic "github.com/cherish-chat/xxim-server/app/gateway/internal/logic/gatewayservice"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/types"
	"github.com/cherish-chat/xxim-server/common/i18n"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"google.golang.org/protobuf/proto"
)

type ConnectionHandler struct {
	svcCtx *svc.ServiceContext
}

func NewConnectionHandler(svcCtx *svc.ServiceContext) *ConnectionHandler {
	return &ConnectionHandler{svcCtx: svcCtx}
}

// VerifyConnection 验证连接
func (h *ConnectionHandler) VerifyConnection(ctx context.Context, connection *gatewayservicelogic.Connection, apiRequest *pb.GatewayApiRequest) (pb.ResponseCode, []byte, error) {
	request := &pb.VerifyConnectionReq{}
	response := &pb.VerifyConnectionResp{}
	err := proto.Unmarshal(apiRequest.Body, request)
	if err != nil {
		responseHeader := i18n.NewInvalidDataError(err.Error())
		response.SetHeader(responseHeader)
		return pb.ResponseCode_INVALID_DATA, types.MarshalWriteData(&pb.GatewayApiResponse{
			RequestId: apiRequest.RequestId,
			Path:      apiRequest.Path,
			Header:    i18n.NewInvalidDataError(err.Error()),
			Body:      utils.Proto.Marshal(response),
		}), err
	}
	response, err = gatewayservicelogic.NewVerifyConnectionLogic(ctx, h.svcCtx).VerifyConnection_(connection, request)
	if err != nil {
		responseHeader := i18n.NewServerError(h.svcCtx.Config.Mode, err)
		response.SetHeader(responseHeader)
		return pb.ResponseCode_SERVER_ERROR, types.MarshalWriteData(&pb.GatewayApiResponse{
			Header:    responseHeader,
			Body:      utils.Proto.Marshal(response),
			RequestId: apiRequest.RequestId,
			Path:      apiRequest.Path,
		}), err
	}
	gatewayservicelogic.ConnectionLogic.OnVerified(connection)
	responseHeader := i18n.NewOkHeader()
	return pb.ResponseCode_SUCCESS, types.MarshalWriteData(&pb.GatewayApiResponse{
		Header:    responseHeader,
		RequestId: apiRequest.RequestId,
		Path:      apiRequest.Path,
		Body:      utils.Proto.Marshal(response),
	}), nil
}

func (h *ConnectionHandler) AuthenticationConnection(ctx context.Context, connection *gatewayservicelogic.Connection, apiRequest *pb.GatewayApiRequest) (pb.ResponseCode, []byte, error) {
	request := &pb.AuthenticationConnectionReq{}
	response := &pb.AuthenticationConnectionResp{}
	err := proto.Unmarshal(apiRequest.Body, request)
	if err != nil {
		responseHeader := i18n.NewInvalidDataError(err.Error())
		response.SetHeader(responseHeader)
		return pb.ResponseCode_INVALID_DATA, types.MarshalWriteData(&pb.GatewayApiResponse{
			RequestId: apiRequest.RequestId,
			Path:      apiRequest.Path,
			Header:    i18n.NewInvalidDataError(err.Error()),
			Body:      utils.Proto.Marshal(response),
		}), err
	}
	response, err = gatewayservicelogic.NewAuthenticationConnectionLogic(ctx, h.svcCtx).AuthenticationConnection_(connection, request)
	if err != nil {
		responseHeader := i18n.NewServerError(h.svcCtx.Config.Mode, err)
		response.SetHeader(responseHeader)
		return pb.ResponseCode_SERVER_ERROR, types.MarshalWriteData(&pb.GatewayApiResponse{
			Header:    responseHeader,
			Body:      utils.Proto.Marshal(response),
			RequestId: apiRequest.RequestId,
			Path:      apiRequest.Path,
		}), err
	}
	responseHeader := i18n.NewOkHeader()
	return pb.ResponseCode_SUCCESS, types.MarshalWriteData(&pb.GatewayApiResponse{
		Header:    responseHeader,
		RequestId: apiRequest.RequestId,
		Path:      apiRequest.Path,
		Body:      utils.Proto.Marshal(response),
	}), nil
}
