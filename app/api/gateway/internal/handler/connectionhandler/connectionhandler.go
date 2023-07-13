package connectionhandler

import (
	"context"
	"github.com/cherish-chat/xxim-proto/peerpb"
	internalservicelogic "github.com/cherish-chat/xxim-server/app/api/gateway/internal/logic/connectionmanager"
	interfaceservicelogic "github.com/cherish-chat/xxim-server/app/api/gateway/internal/logic/interfaceservice"
	"github.com/cherish-chat/xxim-server/app/api/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/app/api/gateway/internal/types"
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
func (h *ConnectionHandler) VerifyConnection(ctx context.Context, connection *internalservicelogic.Connection, apiRequest *peerpb.GatewayApiRequest) (peerpb.ResponseCode, []byte, error) {
	request := &peerpb.VerifyConnectionReq{}
	response := &peerpb.VerifyConnectionResp{}
	err := proto.Unmarshal(apiRequest.Body, request)
	if err != nil {
		responseHeader := peerpb.NewInvalidDataError(err.Error())
		response.SetHeader(responseHeader)
		return peerpb.ResponseCode_INVALID_DATA, types.MarshalWriteData(&peerpb.GatewayApiResponse{
			RequestId: apiRequest.RequestId,
			Path:      apiRequest.Path,
			Header:    peerpb.NewInvalidDataError(err.Error()),
			Body:      utils.Proto.Marshal(response),
		}), err
	}
	response, err = interfaceservicelogic.NewVerifyConnectionLogic(ctx, h.svcCtx).VerifyConnection_(connection, request)
	if err != nil {
		responseHeader := peerpb.NewServerError(h.svcCtx.Config.Mode, err)
		response.SetHeader(responseHeader)
		return peerpb.ResponseCode_SERVER_ERROR, types.MarshalWriteData(&peerpb.GatewayApiResponse{
			Header:    responseHeader,
			Body:      utils.Proto.Marshal(response),
			RequestId: apiRequest.RequestId,
			Path:      apiRequest.Path,
		}), err
	}
	internalservicelogic.ConnectionLogic.OnVerified(connection)
	responseHeader := peerpb.NewOkHeader()
	return peerpb.ResponseCode_SUCCESS, types.MarshalWriteData(&peerpb.GatewayApiResponse{
		Header:    responseHeader,
		RequestId: apiRequest.RequestId,
		Path:      apiRequest.Path,
		Body:      utils.Proto.Marshal(response),
	}), nil
}

func (h *ConnectionHandler) AuthenticationConnection(ctx context.Context, connection *internalservicelogic.Connection, apiRequest *peerpb.GatewayApiRequest) (peerpb.ResponseCode, []byte, error) {
	request := &peerpb.AuthConnectionReq{}
	response := &peerpb.AuthConnectionResp{}
	err := proto.Unmarshal(apiRequest.Body, request)
	if err != nil {
		responseHeader := peerpb.NewInvalidDataError(err.Error())
		response.SetHeader(responseHeader)
		return peerpb.ResponseCode_INVALID_DATA, types.MarshalWriteData(&peerpb.GatewayApiResponse{
			RequestId: apiRequest.RequestId,
			Path:      apiRequest.Path,
			Header:    peerpb.NewInvalidDataError(err.Error()),
			Body:      utils.Proto.Marshal(response),
		}), err
	}
	response, err = interfaceservicelogic.NewAuthConnectionLogic(ctx, h.svcCtx).AuthConnection_(connection, request)
	if err != nil {
		responseHeader := peerpb.NewServerError(h.svcCtx.Config.Mode, err)
		response.SetHeader(responseHeader)
		return peerpb.ResponseCode_SERVER_ERROR, types.MarshalWriteData(&peerpb.GatewayApiResponse{
			Header:    responseHeader,
			Body:      utils.Proto.Marshal(response),
			RequestId: apiRequest.RequestId,
			Path:      apiRequest.Path,
		}), err
	}
	internalservicelogic.ConnectionLogic.OnLogin(connection)
	responseHeader := peerpb.NewOkHeader()
	return peerpb.ResponseCode_SUCCESS, types.MarshalWriteData(&peerpb.GatewayApiResponse{
		Header:    responseHeader,
		RequestId: apiRequest.RequestId,
		Path:      apiRequest.Path,
		Body:      utils.Proto.Marshal(response),
	}), nil
}
