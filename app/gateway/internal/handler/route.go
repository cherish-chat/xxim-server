package handler

import (
	"context"
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

func AddRoute[REQ ReqInterface, RESP RespInterface](path string, route Route[REQ, RESP]) {
	request := route.NewRequest()
	//routeMap[path] = func(ctx context.Context, c *types.UserConn, body IBody) (*pb.GatewayApiResponse, error) {
	//	return UnifiedHandleWebsocket(ctx, path, c, body, request, route.Do)
	//}
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
			buf := MarshalHttp(requestHeader, response)
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
		buf := MarshalHttp(requestHeader, response)
		ctx.Data(200, requestHeader.Encoding.ContentType(), buf)
		return
	}
}

func SetupRoutes(svcCtx *svc.ServiceContext, engine *gin.Engine) {
	AddRoute("/v1/gateway/getUserConnection", Route[*pb.GatewayGetUserConnectionReq, *pb.GatewayGetUserConnectionResp]{
		NewRequest: func() *pb.GatewayGetUserConnectionReq {
			return &pb.GatewayGetUserConnectionReq{}
		},
		Do: svcCtx.GatewayService().GatewayGetUserConnection,
	})
	apiGroup := engine.Group("/api")
	for path, handlerFunc := range httpRouteMap {
		apiGroup.POST(path, handlerFunc)
	}
}
