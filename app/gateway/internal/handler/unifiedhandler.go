package handler

import (
	"context"
	"errors"
	gatewayservicelogic "github.com/cherish-chat/xxim-server/app/gateway/internal/logic/gatewayservice"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/types"
	"github.com/cherish-chat/xxim-server/common/i18n"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

func UnifiedHandleUniversal[REQ ReqInterface, RESP RespInterface](
	svcCtx *svc.ServiceContext,
	ctx context.Context,
	connection *gatewayservicelogic.Connection,
	apiRequest *pb.GatewayApiRequest,
	route Route[REQ, RESP],
) (REQ, RESP, *pb.GatewayApiResponse, error) {
	request := route.RequestPool.NewRequest()
	response := route.ResponsePool.NewResponse()
	do := route.Do
	// 请求体中的数据 反序列化到 request 中
	// 判断是json还是protobuf
	requestHeader := apiRequest.Header
	err := proto.Unmarshal(apiRequest.Body, request)
	if err != nil {
		responseHeader := i18n.NewInvalidDataError(err.Error())
		response.SetHeader(responseHeader)
		return request, response, &pb.GatewayApiResponse{
			RequestId: apiRequest.RequestId,
			Path:      apiRequest.Path,
			Header:    i18n.NewInvalidDataError(err.Error()),
			Body:      utils.Proto.Marshal(response),
		}, err
	}
	request.SetHeader(requestHeader)
	if requestHeader == nil {
		responseHeader := i18n.NewInvalidDataError("invalid request header")
		response.SetHeader(responseHeader)
		return request, response, &pb.GatewayApiResponse{
			Header:    responseHeader,
			Body:      utils.Proto.Marshal(response),
			RequestId: apiRequest.RequestId,
			Path:      apiRequest.Path,
		}, nil
	}

	//beforeRequest
	userBeforeRequestResp, err := svcCtx.CallbackService.UserBeforeRequest(ctx, &pb.UserBeforeRequestReq{
		Header: requestHeader,
		Path:   apiRequest.Path,
	})
	if err != nil {
		responseHeader := i18n.NewServerError(svcCtx.Config.Mode, err)
		response.SetHeader(responseHeader)
		return request, response, &pb.GatewayApiResponse{
			Header:    responseHeader,
			RequestId: apiRequest.RequestId,
			Path:      apiRequest.Path,
			Body:      utils.Proto.Marshal(response),
		}, err
	}

	var result *pb.GatewayApiResponse
	body, _ := proto.Marshal(userBeforeRequestResp)
	if len(body) > 0 {
		responseHeader := userBeforeRequestResp.GetHeader()
		if responseHeader != nil && responseHeader.Code != pb.ResponseCode_SUCCESS {
			return request, response, &pb.GatewayApiResponse{
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
			responseHeader = i18n.NewOkHeader()
		}
		result = &pb.GatewayApiResponse{
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
			responseHeader := i18n.NewServerError(svcCtx.Config.Mode, err)
			response.SetHeader(responseHeader)
			result = &pb.GatewayApiResponse{
				Header:    responseHeader,
				RequestId: apiRequest.RequestId,
				Path:      apiRequest.Path,
				Body:      utils.Proto.Marshal(response),
			}
		} else {
			responseHeader := i18n.NewOkHeader()
			response.SetHeader(responseHeader)
			result = &pb.GatewayApiResponse{
				Header:    responseHeader,
				RequestId: apiRequest.RequestId,
				Path:      apiRequest.Path,
				Body:      utils.Proto.Marshal(response),
			}
		}
	}
	return request, response, result, err
}
