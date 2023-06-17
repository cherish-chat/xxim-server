package handler

import (
	"context"
	gatewayservicelogic "github.com/cherish-chat/xxim-server/app/gateway/internal/logic/gatewayservice"
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
var wsRouteMap = map[string]func(ctx context.Context, connection *gatewayservicelogic.WsConnection, c *pb.GatewayApiRequest) (pb.ResponseCode, []byte, error){}

func AddUnifiedRoute[REQ ReqInterface, RESP RespInterface](svcCtx *svc.ServiceContext, path string, route Route[REQ, RESP]) {
	request := route.NewRequest()
	// http
	AddHttpRoute(path, func(ctx *gin.Context) {
		var response *pb.GatewayApiResponse
		var err error
		response, err = UnifiedHandleHttp(svcCtx, ctx, request, route.Do)
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
	})
	// ws
	AddWsRoute(svcCtx, path, func(ctx context.Context, connection *gatewayservicelogic.WsConnection, apiRequest *pb.GatewayApiRequest) (pb.ResponseCode, []byte, error) {
		var response *pb.GatewayApiResponse
		var err error
		requestHeader := connection.Header
		response, err = UnifiedHandleWs(svcCtx, ctx, connection, apiRequest, request, route.Do)
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
	})
}

func AddHttpRoute(path string, handlerFunc gin.HandlerFunc) {
	httpRouteMap[path] = handlerFunc
}

func AddWsRoute(svcCtx *svc.ServiceContext, path string, handlerFunc func(ctx context.Context, connection *gatewayservicelogic.WsConnection, c *pb.GatewayApiRequest) (pb.ResponseCode, []byte, error)) {
	wsRouteMap[path] = handlerFunc
}

func SetupRoutes(svcCtx *svc.ServiceContext, engine *gin.Engine) {
	// gateway api
	{
		// GatewayGetUserConnectionReq GatewayGetUserConnectionResp
		AddUnifiedRoute(svcCtx, "/v1/gateway/getUserConnection", Route[*pb.GatewayGetUserConnectionReq, *pb.GatewayGetUserConnectionResp]{
			NewRequest: func() *pb.GatewayGetUserConnectionReq {
				return &pb.GatewayGetUserConnectionReq{}
			},
			Do: svcCtx.GatewayService().GatewayGetUserConnection,
		})
		// GatewayBatchGetUserConnectionReq GatewayBatchGetUserConnectionResp
		AddUnifiedRoute(svcCtx, "/v1/gateway/batchGetUserConnection", Route[*pb.GatewayBatchGetUserConnectionReq, *pb.GatewayBatchGetUserConnectionResp]{
			NewRequest: func() *pb.GatewayBatchGetUserConnectionReq {
				return &pb.GatewayBatchGetUserConnectionReq{}
			},
			Do: svcCtx.GatewayService().GatewayBatchGetUserConnection,
		})
		// GatewayGetConnectionByFilterReq GatewayGetConnectionByFilterResp
		AddUnifiedRoute(svcCtx, "/v1/gateway/getConnectionByFilter", Route[*pb.GatewayGetConnectionByFilterReq, *pb.GatewayGetConnectionByFilterResp]{
			NewRequest: func() *pb.GatewayGetConnectionByFilterReq {
				return &pb.GatewayGetConnectionByFilterReq{}
			},
			Do: svcCtx.GatewayService().GatewayGetConnectionByFilter,
		})
		// GatewayWriteDataToWsReq GatewayWriteDataToWsResp
		AddUnifiedRoute(svcCtx, "/v1/gateway/writeDataToWs", Route[*pb.GatewayWriteDataToWsReq, *pb.GatewayWriteDataToWsResp]{
			NewRequest: func() *pb.GatewayWriteDataToWsReq {
				return &pb.GatewayWriteDataToWsReq{}
			},
			Do: svcCtx.GatewayService().GatewayWriteDataToWs,
		})
		// GatewayKickWsReq GatewayKickWsResp
		AddUnifiedRoute(svcCtx, "/v1/gateway/kickWs", Route[*pb.GatewayKickWsReq, *pb.GatewayKickWsResp]{
			NewRequest: func() *pb.GatewayKickWsReq {
				return &pb.GatewayKickWsReq{}
			},
			Do: svcCtx.GatewayService().GatewayKickWs,
		})
		// GatewayKeepAliveReq GatewayKeepAliveResp
		AddWsRoute(svcCtx, "/v1/gateway/white/keepAlive", KeepAliveHandler(svcCtx))
	}
	// user api
	{
		// UserRegisterReq UserRegisterResp
		AddUnifiedRoute(svcCtx, "/v1/user/white/userRegister", Route[*pb.UserRegisterReq, *pb.UserRegisterResp]{
			NewRequest: func() *pb.UserRegisterReq {
				return &pb.UserRegisterReq{}
			},
			Do: svcCtx.AccountService.UserRegister,
		})
		// UserAccessTokenReq UserAccessTokenResp
		AddUnifiedRoute(svcCtx, "/v1/user/white/userAccessToken", Route[*pb.UserAccessTokenReq, *pb.UserAccessTokenResp]{
			NewRequest: func() *pb.UserAccessTokenReq {
				return &pb.UserAccessTokenReq{}
			},
			Do: svcCtx.AccountService.UserAccessToken,
		})
		// CreateRobotReq CreateRobotResp
		AddUnifiedRoute(svcCtx, "/v1/user/createRobot", Route[*pb.CreateRobotReq, *pb.CreateRobotResp]{
			NewRequest: func() *pb.CreateRobotReq {
				return &pb.CreateRobotReq{}
			},
			Do: svcCtx.AccountService.CreateRobot,
		})
		// RefreshUserAccessTokenReq RefreshUserAccessTokenResp
		AddUnifiedRoute(svcCtx, "/v1/user/refreshUserAccessToken", Route[*pb.RefreshUserAccessTokenReq, *pb.RefreshUserAccessTokenResp]{
			NewRequest: func() *pb.RefreshUserAccessTokenReq {
				return &pb.RefreshUserAccessTokenReq{}
			},
			Do: svcCtx.AccountService.RefreshUserAccessToken,
		})
		// RevokeUserAccessTokenReq RevokeUserAccessTokenResp
		AddUnifiedRoute(svcCtx, "/v1/user/revokeUserAccessToken", Route[*pb.RevokeUserAccessTokenReq, *pb.RevokeUserAccessTokenResp]{
			NewRequest: func() *pb.RevokeUserAccessTokenReq {
				return &pb.RevokeUserAccessTokenReq{}
			},
			Do: svcCtx.AccountService.RevokeUserAccessToken,
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
