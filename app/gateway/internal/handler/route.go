package handler

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/logic"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/common/i18n"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

type ReqInterface interface {
	proto.Message
	SetHeader(*pb.RequestHeader)
	GetHeader() *pb.RequestHeader
}

type RespInterface interface {
	proto.Message
	GetHeader() *pb.ResponseHeader
}

type Route[REQ ReqInterface, RESP RespInterface] struct {
	NewRequest func() REQ
	Do         func(ctx context.Context, req REQ, opts ...grpc.CallOption) (RESP, error)
}

//var routeMap = map[string]func(ctx context.Context, c *types.UserConn, body IBody) (*pb.ResponseBody, error){}

var httpRouteMap = map[string]gin.HandlerFunc{}
var wsRouteMap = map[string]func(ctx context.Context, connection *logic.WsConnection, c *pb.GatewayApiRequest) (pb.ResponseCode, []byte, error){}

func AddUnifiedRoute[REQ ReqInterface, RESP RespInterface](path string, route Route[REQ, RESP]) {
	request := route.NewRequest()
	// http
	httpRouteMap[path] = func(ctx *gin.Context) {
		var response *pb.GatewayApiResponse
		var err error
		response, err = UnifiedHandleHttp(ctx, request, route.Do)
		requestHeader := request.GetHeader()
		if err != nil {
			logx.WithContext(ctx.Request.Context()).Errorf("UnifiedHandleHttp: %s, error: %v", path, err)
			if response == nil {
				response = &pb.GatewayApiResponse{
					Header: i18n.NewServerError(requestHeader),
					Body:   nil,
				}
			}
			buf := MarshalResponse(requestHeader, response)
			ctx.Data(200, requestHeader.Encoding.ContentType(), buf)
			return
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
		buf := MarshalResponse(requestHeader, response)
		ctx.Data(200, requestHeader.Encoding.ContentType(), buf)
		return
	}
	// ws
	wsRouteMap[path] = func(ctx context.Context, connection *logic.WsConnection, apiRequest *pb.GatewayApiRequest) (pb.ResponseCode, []byte, error) {
		var response *pb.GatewayApiResponse
		var err error
		requestHeader := connection.Header
		response, err = UnifiedHandleWs(ctx, connection, apiRequest, request, route.Do)
		if err != nil {
			logx.WithContext(ctx).Errorf("UnifiedHandleWs: %s, error: %v", path, err)
			if response == nil {
				response = &pb.GatewayApiResponse{
					Header: i18n.NewServerError(requestHeader),
					Body:   nil,
				}
			}
			return response.Header.Code, MarshalResponse(requestHeader, response), nil
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
		return response.Header.Code, MarshalResponse(requestHeader, response), nil
	}
}

func AddHttpRoute(path string, handlerFunc gin.HandlerFunc) {
	httpRouteMap[path] = handlerFunc
}

func AddWsRoute(path string, handlerFunc func(ctx context.Context, connection *logic.WsConnection, c *pb.GatewayApiRequest) (pb.ResponseCode, []byte, error)) {
	wsRouteMap[path] = handlerFunc
}

func SetupRoutes(svcCtx *svc.ServiceContext, engine *gin.Engine) {
	// gateway api
	{
		AddUnifiedRoute("/v1/gateway/getUserConnection", Route[*pb.GatewayGetUserConnectionReq, *pb.GatewayGetUserConnectionResp]{
			NewRequest: func() *pb.GatewayGetUserConnectionReq {
				return &pb.GatewayGetUserConnectionReq{}
			},
			Do: svcCtx.GatewayService().GatewayGetUserConnection,
		})
	}
	// http
	{
		apiGroup := engine.Group("/api")
		for path, handlerFunc := range httpRouteMap {
			apiGroup.POST(path, handlerFunc)
		}
	}
	// ws
	{
		wsHandler := NewWsHandler(svcCtx)
		engine.GET("/ws", wsHandler.Upgrade)
	}
}
