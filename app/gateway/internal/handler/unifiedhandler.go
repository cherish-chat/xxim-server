package handler

import (
	"context"
	"errors"
	gatewayservicelogic "github.com/cherish-chat/xxim-server/app/gateway/internal/logic/gatewayservice"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/common/i18n"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"strings"
)

func UnifiedHandleHttp[REQ ReqInterface, RESP RespInterface](
	svcCtx *svc.ServiceContext,
	ctx *gin.Context,
	route Route[REQ, RESP],
) (REQ, *pb.GatewayApiResponse, error) {
	request := route.RequestPool.NewRequest()
	response := route.ResponsePool.NewResponse()
	do := route.Do
	// 请求体中的数据 反序列化到 request 中
	contentType := ctx.ContentType()
	encoding := pb.EncodingProto_PROTOBUF
	// 判断是json还是protobuf
	apiRequest := &pb.GatewayApiRequest{}
	if strings.Contains(contentType, "application/json") {
		// json
		err := ctx.ShouldBindJSON(apiRequest)
		if err != nil {
			responseHeader := i18n.NewInvalidDataError(err.Error())
			response.SetHeader(responseHeader)
			return request, &pb.GatewayApiResponse{
				Header: responseHeader,
				Body:   utils.Proto.Marshal(response),
			}, err
		}
		err = utils.Json.Unmarshal(apiRequest.Body, request)
		if err != nil {
			responseHeader := i18n.NewInvalidDataError(err.Error())
			response.SetHeader(responseHeader)
			return request, &pb.GatewayApiResponse{
				RequestId: apiRequest.RequestId,
				Path:      apiRequest.Path,
				Header:    i18n.NewInvalidDataError(err.Error()),
				Body:      utils.Proto.Marshal(response),
			}, err
		}
		request.SetHeader(apiRequest.Header)
		encoding = pb.EncodingProto_JSON
	} else if strings.Contains(contentType, "application/x-protobuf") {
		// protobuf
		body, _ := ctx.GetRawData()
		err := proto.Unmarshal(body, apiRequest)
		if err != nil {
			responseHeader := i18n.NewInvalidDataError(err.Error())
			response.SetHeader(responseHeader)
			return request, &pb.GatewayApiResponse{
				RequestId: apiRequest.RequestId,
				Path:      apiRequest.Path,
				Header:    i18n.NewInvalidDataError(err.Error()),
				Body:      utils.Proto.Marshal(response),
			}, err
		}
		err = proto.Unmarshal(apiRequest.Body, request)
		if err != nil {
			responseHeader := i18n.NewInvalidDataError(err.Error())
			response.SetHeader(responseHeader)
			return request, &pb.GatewayApiResponse{
				RequestId: apiRequest.RequestId,
				Path:      apiRequest.Path,
				Header:    i18n.NewInvalidDataError(err.Error()),
				Body:      utils.Proto.Marshal(response),
			}, err
		}
		request.SetHeader(apiRequest.Header)
	} else {
		responseHeader := i18n.NewInvalidDataError("invalid content type, please use application/json or application/x-protobuf")
		response.SetHeader(responseHeader)
		return request, &pb.GatewayApiResponse{
			RequestId: apiRequest.RequestId,
			Path:      apiRequest.Path,
			Header:    responseHeader,
			Body:      utils.Proto.Marshal(response),
		}, nil
	}
	requestHeader := request.GetHeader()
	if requestHeader == nil {
		responseHeader := i18n.NewInvalidDataError("invalid request header")
		response.SetHeader(responseHeader)
		return request, &pb.GatewayApiResponse{
			RequestId: apiRequest.RequestId,
			Path:      apiRequest.Path,
			Header:    responseHeader,
			Body:      utils.Proto.Marshal(response),
		}, nil
	}
	requestHeader.ClientIp = utils.Http.GetClientIP(ctx.Request)
	requestHeader.Encoding = encoding
	requestHeader.GatewayPodIp = utils.GetPodIp()

	// beforeRequest
	userBeforeRequestResp, err := svcCtx.CallbackService.UserBeforeRequest(ctx.Request.Context(), &pb.UserBeforeRequestReq{
		Header: requestHeader,
		Path:   apiRequest.Path,
	})
	if err != nil {
		responseHeader := i18n.NewServerError(requestHeader, svcCtx.Config.Mode, err)
		response.SetHeader(responseHeader)
		return request, &pb.GatewayApiResponse{
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
			return request, &pb.GatewayApiResponse{
				Header:    responseHeader,
				RequestId: apiRequest.RequestId,
				Path:      apiRequest.Path,
				Body:      MarshalResponse(requestHeader, userBeforeRequestResp),
			}, nil
		}
	}

	requestHeader.UserId = userBeforeRequestResp.UserId

	resp, err := do(ctx.Request.Context(), request)
	body, _ = proto.Marshal(resp)
	if len(body) > 0 {
		responseHeader := resp.GetHeader()
		if responseHeader == nil {
			responseHeader = i18n.NewOkHeader()
		}
		result = &pb.GatewayApiResponse{
			Header:    responseHeader,
			RequestId: apiRequest.RequestId,
			Path:      apiRequest.Path,
			Body:      MarshalResponse(requestHeader, resp),
		}
	} else {
		if err != nil {
			statusErr, ok := status.FromError(err)
			if ok {
				err = errors.New(statusErr.Message())
			}
			responseHeader := i18n.NewServerError(requestHeader, svcCtx.Config.Mode, err)
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
	return request, result, err
}

func UnifiedHandleUniversal[REQ ReqInterface, RESP RespInterface](
	svcCtx *svc.ServiceContext,
	ctx context.Context,
	connection *gatewayservicelogic.UniversalConnection,
	apiRequest *pb.GatewayApiRequest,
	route Route[REQ, RESP],
) (REQ, *pb.GatewayApiResponse, error) {
	request := route.RequestPool.NewRequest()
	response := route.ResponsePool.NewResponse()
	do := route.Do
	// 请求体中的数据 反序列化到 request 中
	// 判断是json还是protobuf
	requestHeader := apiRequest.Header
	if requestHeader.Encoding == pb.EncodingProto_JSON {
		err := utils.Json.Unmarshal(apiRequest.Body, request)
		if err != nil {
			responseHeader := i18n.NewInvalidDataError(err.Error())
			response.SetHeader(responseHeader)
			return request, &pb.GatewayApiResponse{
				RequestId: apiRequest.RequestId,
				Path:      apiRequest.Path,
				Header:    responseHeader,
				Body:      utils.Proto.Marshal(response),
			}, err
		}
	} else if requestHeader.Encoding == pb.EncodingProto_PROTOBUF {
		err := proto.Unmarshal(apiRequest.Body, request)
		if err != nil {
			responseHeader := i18n.NewInvalidDataError(err.Error())
			response.SetHeader(responseHeader)
			return request, &pb.GatewayApiResponse{
				RequestId: apiRequest.RequestId,
				Path:      apiRequest.Path,
				Header:    i18n.NewInvalidDataError(err.Error()),
				Body:      utils.Proto.Marshal(response),
			}, err
		}
	} else {
		responseHeader := i18n.NewInvalidDataError("invalid content type, please use application/json or application/x-protobuf")
		response.SetHeader(responseHeader)
		return request, &pb.GatewayApiResponse{
			Header: responseHeader,
			Body:   utils.Proto.Marshal(response),
		}, nil
	}
	request.SetHeader(requestHeader)
	if requestHeader == nil {
		responseHeader := i18n.NewInvalidDataError("invalid request header")
		response.SetHeader(responseHeader)
		return request, &pb.GatewayApiResponse{
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
		responseHeader := i18n.NewServerError(requestHeader, svcCtx.Config.Mode, err)
		response.SetHeader(responseHeader)
		return request, &pb.GatewayApiResponse{
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
			return request, &pb.GatewayApiResponse{
				Header:    responseHeader,
				RequestId: apiRequest.RequestId,
				Path:      apiRequest.Path,
				Body:      MarshalResponse(requestHeader, userBeforeRequestResp),
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
			Body:      MarshalResponse(requestHeader, resp),
		}
	} else {
		if err != nil {
			statusErr, ok := status.FromError(err)
			if ok {
				err = errors.New(statusErr.Message())
			}
			responseHeader := i18n.NewServerError(requestHeader, svcCtx.Config.Mode, err)
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
	return request, result, err
}

func MarshalResponse(requestHeader *pb.RequestHeader, data proto.Message) []byte {
	if requestHeader == nil {
		return nil
	}
	switch requestHeader.Encoding {
	case pb.EncodingProto_PROTOBUF:
		protobuf, _ := proto.Marshal(data)
		return protobuf
	case pb.EncodingProto_JSON:
		json, _ := utils.Json.Marshal(data)
		return json
	default:
		return nil
	}
}

func MarshalWriteData(requestHeader *pb.RequestHeader, data *pb.GatewayApiResponse) []byte {
	if requestHeader == nil {
		return nil
	}
	writeData := &pb.GatewayWriteDataContent{
		DataType: pb.GatewayWriteDataType_Response,
		Response: data,
		Message:  nil,
		Notice:   nil,
	}
	switch requestHeader.Encoding {
	case pb.EncodingProto_PROTOBUF:
		protobuf, _ := proto.Marshal(writeData)
		return protobuf
	case pb.EncodingProto_JSON:
		json, _ := utils.Json.Marshal(writeData)
		return json
	default:
		return nil
	}
}
