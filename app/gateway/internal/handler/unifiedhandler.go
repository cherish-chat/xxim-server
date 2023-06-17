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
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
	"strings"
)

func UnifiedHandleHttp[REQ ReqInterface, RESP RespInterface](
	svcCtx *svc.ServiceContext,
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

	// beforeRequest
	userBeforeRequestResp, err := svcCtx.CallbackService.UserBeforeRequest(ctx.Request.Context(), &pb.UserBeforeRequestReq{
		Header: requestHeader,
		Path:   apiRequest.Path,
	})
	if err != nil {
		return &pb.GatewayApiResponse{
			Header:    i18n.NewServerError(requestHeader, svcCtx.Config.Mode, err),
			RequestId: apiRequest.RequestId,
			Path:      apiRequest.Path,
			Body:      nil,
		}, err
	}
	requestHeader.UserId = userBeforeRequestResp.UserId

	var result *pb.GatewayApiResponse
	body, _ := proto.Marshal(userBeforeRequestResp)
	if len(body) > 0 {
		responseHeader := userBeforeRequestResp.GetHeader()
		if responseHeader != nil && responseHeader.Code != pb.ResponseCode_SUCCESS {
			return &pb.GatewayApiResponse{
				Header:    responseHeader,
				RequestId: apiRequest.RequestId,
				Path:      apiRequest.Path,
				Body:      MarshalResponse(requestHeader, userBeforeRequestResp),
			}, nil
		}
	}

	response, err := do(ctx.Request.Context(), request)
	body, _ = proto.Marshal(response)
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
		if err != nil {
			statusErr, ok := status.FromError(err)
			if ok {
				err = errors.New(statusErr.Message())
			}
			result = &pb.GatewayApiResponse{
				Header:    i18n.NewServerError(requestHeader, svcCtx.Config.Mode, err),
				RequestId: apiRequest.RequestId,
				Path:      apiRequest.Path,
				Body:      nil,
			}
		} else {
			result = &pb.GatewayApiResponse{
				Header:    i18n.NewOkHeader(),
				RequestId: apiRequest.RequestId,
				Path:      apiRequest.Path,
				Body:      nil,
			}
		}
	}
	return result, err
}

func UnifiedHandleWs[REQ ReqInterface, RESP RespInterface](
	svcCtx *svc.ServiceContext,
	ctx context.Context,
	connection *gatewayservicelogic.WsConnection,
	apiRequest *pb.GatewayApiRequest,
	request REQ,
	do func(ctx context.Context, req REQ, opts ...grpc.CallOption) (RESP, error),
) (*pb.GatewayApiResponse, error) {
	// 请求体中的数据 反序列化到 request 中
	// 判断是json还是protobuf
	requestHeader := apiRequest.Header
	if requestHeader.Encoding == pb.EncodingProto_JSON {
		err := utils.Json.Unmarshal(apiRequest.Body, request)
		if err != nil {
			return &pb.GatewayApiResponse{
				RequestId: apiRequest.RequestId,
				Path:      apiRequest.Path,
				Header:    i18n.NewInvalidDataError(err.Error()),
				Body:      nil,
			}, err
		}
	} else if requestHeader.Encoding == pb.EncodingProto_PROTOBUF {
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
	request.SetHeader(requestHeader)
	if requestHeader == nil {
		return &pb.GatewayApiResponse{
			Header:    i18n.NewInvalidDataError("invalid request header"),
			Body:      nil,
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
		return &pb.GatewayApiResponse{
			Header:    i18n.NewServerError(requestHeader, svcCtx.Config.Mode, err),
			RequestId: apiRequest.RequestId,
			Path:      apiRequest.Path,
			Body:      nil,
		}, err
	}
	requestHeader.UserId = userBeforeRequestResp.UserId

	var result *pb.GatewayApiResponse
	body, _ := proto.Marshal(userBeforeRequestResp)
	if len(body) > 0 {
		responseHeader := userBeforeRequestResp.GetHeader()
		if responseHeader != nil && responseHeader.Code != pb.ResponseCode_SUCCESS {
			return &pb.GatewayApiResponse{
				Header:    responseHeader,
				RequestId: apiRequest.RequestId,
				Path:      apiRequest.Path,
				Body:      MarshalResponse(requestHeader, userBeforeRequestResp),
			}, nil
		}
	}

	response, err := do(ctx, request)
	body, _ = proto.Marshal(response)
	if len(body) > 0 {
		responseHeader := response.GetHeader()
		if responseHeader == nil {
			responseHeader = i18n.NewOkHeader()
		}
		result = &pb.GatewayApiResponse{
			RequestId: apiRequest.RequestId,
			Path:      apiRequest.Path,
			Header:    responseHeader,
			Body:      MarshalResponse(requestHeader, response),
		}
	} else {
		if err != nil {
			statusErr, ok := status.FromError(err)
			if ok {
				err = errors.New(statusErr.Message())
			}
			result = &pb.GatewayApiResponse{
				Header:    i18n.NewServerError(requestHeader, svcCtx.Config.Mode, err),
				RequestId: apiRequest.RequestId,
				Path:      apiRequest.Path,
				Body:      nil,
			}
		} else {
			result = &pb.GatewayApiResponse{
				Header:    i18n.NewOkHeader(),
				RequestId: apiRequest.RequestId,
				Path:      apiRequest.Path,
				Body:      nil,
			}
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
