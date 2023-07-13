package handler

import (
	"context"
	"errors"
	"github.com/cherish-chat/xxim-proto/peerpb"
	internalservicelogic "github.com/cherish-chat/xxim-server/app/api/gateway/internal/logic/connectionmanager"
	"github.com/cherish-chat/xxim-server/app/api/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/app/api/gateway/internal/types"
	"github.com/cherish-chat/xxim-server/common/utils"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func UnifiedHandleUniversal[REQ ReqInterface, RESP RespInterface](
	svcCtx *svc.ServiceContext,
	ctx context.Context,
	connection *internalservicelogic.Connection,
	apiRequest *peerpb.GatewayApiRequest,
	route Route[REQ, RESP],
) (REQ, RESP, *peerpb.GatewayApiResponse, error) {
	request := route.RequestPool.NewRequest()
	response := route.ResponsePool.NewResponse()
	do := route.Do
	// 请求体中的数据 反序列化到 request 中
	// 判断是json还是protobuf
	requestHeader := apiRequest.Header
	err := proto.Unmarshal(apiRequest.Body, request)
	if err != nil {
		responseHeader := peerpb.NewInvalidDataError(err.Error())
		response.SetHeader(responseHeader)
		return request, response, &peerpb.GatewayApiResponse{
			RequestId: apiRequest.RequestId,
			Path:      apiRequest.Path,
			Header:    peerpb.NewInvalidDataError(err.Error()),
			Body:      utils.Proto.Marshal(response),
		}, err
	}
	request.SetHeader(requestHeader)
	if requestHeader == nil {
		responseHeader := peerpb.NewInvalidDataError("invalid request header")
		response.SetHeader(responseHeader)
		return request, response, &peerpb.GatewayApiResponse{
			Header:    responseHeader,
			Body:      utils.Proto.Marshal(response),
			RequestId: apiRequest.RequestId,
			Path:      apiRequest.Path,
		}, nil
	}

	//beforeRequest
	userBeforeRequestResp, err := svcCtx.CallbackService.UserBeforeRequest(ctx, &peerpb.UserBeforeRequestReq{
		Header: requestHeader,
		Path:   apiRequest.Path,
	})
	if err != nil {
		responseHeader := peerpb.NewServerError(svcCtx.Config.Mode, err)
		response.SetHeader(responseHeader)
		return request, response, &peerpb.GatewayApiResponse{
			Header:    responseHeader,
			RequestId: apiRequest.RequestId,
			Path:      apiRequest.Path,
			Body:      utils.Proto.Marshal(response),
		}, err
	}

	var result *peerpb.GatewayApiResponse
	body, _ := proto.Marshal(userBeforeRequestResp)
	if len(body) > 0 {
		responseHeader := userBeforeRequestResp.GetHeader()
		if responseHeader != nil && responseHeader.Code != peerpb.ResponseCode_SUCCESS {
			return request, response, &peerpb.GatewayApiResponse{
				Header:    responseHeader,
				RequestId: apiRequest.RequestId,
				Path:      apiRequest.Path,
				Body:      types.MarshalResponse(userBeforeRequestResp),
			}, nil
		}
	}
	resp, err := do(ctx, request)
	body, _ = proto.Marshal(resp)
	if len(body) > 0 {
		responseHeader := resp.GetHeader()
		if responseHeader == nil {
			responseHeader = peerpb.NewOkHeader()
		}
		result = &peerpb.GatewayApiResponse{
			RequestId: apiRequest.RequestId,
			Path:      apiRequest.Path,
			Header:    responseHeader,
			Body:      types.MarshalResponse(resp),
		}
	} else {
		if err != nil {
			statusErr, ok := status.FromError(err)
			if ok {
				err = errors.New(statusErr.Message())
			}
			responseHeader := peerpb.NewServerError(svcCtx.Config.Mode, err)
			response.SetHeader(responseHeader)
			result = &peerpb.GatewayApiResponse{
				Header:    responseHeader,
				RequestId: apiRequest.RequestId,
				Path:      apiRequest.Path,
				Body:      utils.Proto.Marshal(response),
			}
		} else {
			responseHeader := peerpb.NewOkHeader()
			response.SetHeader(responseHeader)
			result = &peerpb.GatewayApiResponse{
				Header:    responseHeader,
				RequestId: apiRequest.RequestId,
				Path:      apiRequest.Path,
				Body:      utils.Proto.Marshal(response),
			}
		}
	}
	return request, response, result, err
}
