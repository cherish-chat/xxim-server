package handler

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/logic"
	"github.com/cherish-chat/xxim-server/common/i18n"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"strings"
)

func UnifiedHandleHttp[REQ ReqInterface, RESP RespInterface](
	ctx *gin.Context,
	request REQ,
	do func(ctx context.Context, req REQ, opts ...grpc.CallOption) (RESP, error),
) (*pb.GatewayApiResponse, error) {
	// 请求体中的数据 反序列化到 request 中
	contentType := ctx.ContentType()
	encoding := pb.EncodingProto_PROTOBUF
	// 判断是json还是protobuf
	apiRequest := &pb.GatewayApiRequest{}
	if strings.Contains(contentType, "application/json") {
		// json
		err := ctx.ShouldBindJSON(apiRequest)
		if err != nil {
			return &pb.GatewayApiResponse{
				Header: i18n.NewInvalidDataError(err.Error()),
				Body:   nil,
			}, err
		}
		err = utils.Json.Unmarshal(apiRequest.Body, request)
		if err != nil {
			return &pb.GatewayApiResponse{
				RequestId: apiRequest.RequestId,
				Path:      apiRequest.Path,
				Header:    i18n.NewInvalidDataError(err.Error()),
				Body:      nil,
			}, err
		}
		request.SetHeader(apiRequest.Header)
		encoding = pb.EncodingProto_JSON
	} else if strings.Contains(contentType, "application/x-protobuf") {
		// protobuf
		body, _ := ctx.GetRawData()
		apiRequest := &pb.GatewayApiRequest{}
		err := proto.Unmarshal(body, apiRequest)
		if err != nil {
			return &pb.GatewayApiResponse{
				RequestId: apiRequest.RequestId,
				Path:      apiRequest.Path,
				Header:    i18n.NewInvalidDataError(err.Error()),
				Body:      nil,
			}, err
		}
		err = proto.Unmarshal(apiRequest.Body, request)
		if err != nil {
			return &pb.GatewayApiResponse{
				RequestId: apiRequest.RequestId,
				Path:      apiRequest.Path,
				Header:    i18n.NewInvalidDataError(err.Error()),
				Body:      nil,
			}, err
		}
		request.SetHeader(apiRequest.Header)
	} else {
		return &pb.GatewayApiResponse{
			RequestId: apiRequest.RequestId,
			Path:      apiRequest.Path,
			Header:    i18n.NewInvalidDataError("invalid content type, please use application/json or application/x-protobuf"),
			Body:      nil,
		}, nil
	}
	requestHeader := request.GetHeader()
	if requestHeader == nil {
		return &pb.GatewayApiResponse{
			RequestId: apiRequest.RequestId,
			Path:      apiRequest.Path,
			Header:    i18n.NewInvalidDataError("invalid request header"),
			Body:      nil,
		}, nil
	}
	requestHeader.ClientIp = utils.Http.GetClientIP(ctx.Request)
	requestHeader.Encoding = encoding
	requestHeader.GatewayPodIp = utils.GetPodIp()
	response, err := do(ctx.Request.Context(), request)
	var result *pb.GatewayApiResponse
	body, _ := proto.Marshal(response)
	if len(body) > 0 {
		responseHeader := response.GetHeader()
		if responseHeader == nil {
			responseHeader = i18n.NewOkHeader()
		}
		result = &pb.GatewayApiResponse{
			Header:    responseHeader,
			RequestId: apiRequest.RequestId,
			Path:      apiRequest.Path,
			Body:      MarshalResponse(requestHeader, response),
		}
	} else {
		result = &pb.GatewayApiResponse{
			Header:    i18n.NewOkHeader(),
			RequestId: apiRequest.RequestId,
			Path:      apiRequest.Path,
			Body:      nil,
		}
	}
	return result, err
}

func UnifiedHandleWs[REQ ReqInterface, RESP RespInterface](
	ctx context.Context,
	connection *logic.WsConnection,
	apiRequest *pb.GatewayApiRequest,
	request REQ,
	do func(ctx context.Context, req REQ, opts ...grpc.CallOption) (RESP, error),
) (*pb.GatewayApiResponse, error) {
	// 请求体中的数据 反序列化到 request 中
	// 判断是json还是protobuf
	header := apiRequest.Header
	if header.Encoding == pb.EncodingProto_JSON {
		err := utils.Json.Unmarshal(apiRequest.Body, request)
		if err != nil {
			return &pb.GatewayApiResponse{
				RequestId: apiRequest.RequestId,
				Path:      apiRequest.Path,
				Header:    i18n.NewInvalidDataError(err.Error()),
				Body:      nil,
			}, err
		}
	} else if header.Encoding == pb.EncodingProto_PROTOBUF {
		err := proto.Unmarshal(apiRequest.Body, request)
		if err != nil {
			return &pb.GatewayApiResponse{
				RequestId: apiRequest.RequestId,
				Path:      apiRequest.Path,
				Header:    i18n.NewInvalidDataError(err.Error()),
				Body:      nil,
			}, err
		}
	} else {
		return &pb.GatewayApiResponse{
			Header: i18n.NewInvalidDataError("invalid content type, please use application/json or application/x-protobuf"),
			Body:   nil,
		}, nil
	}
	request.SetHeader(header)
	if header == nil {
		return &pb.GatewayApiResponse{
			Header:    i18n.NewInvalidDataError("invalid request header"),
			Body:      nil,
			RequestId: apiRequest.RequestId,
			Path:      apiRequest.Path,
		}, nil
	}
	response, err := do(ctx, request)
	var result *pb.GatewayApiResponse
	body, _ := proto.Marshal(response)
	if len(body) > 0 {
		responseHeader := response.GetHeader()
		if responseHeader == nil {
			responseHeader = i18n.NewOkHeader()
		}
		result = &pb.GatewayApiResponse{
			RequestId: apiRequest.RequestId,
			Path:      apiRequest.Path,
			Header:    responseHeader,
			Body:      MarshalResponse(header, response),
		}
	} else {
		result = &pb.GatewayApiResponse{
			RequestId: apiRequest.RequestId,
			Path:      apiRequest.Path,
			Header:    i18n.NewOkHeader(),
			Body:      nil,
		}
	}
	return result, err
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
