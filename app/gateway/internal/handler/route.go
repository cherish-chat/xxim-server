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
	"sync"
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

// RequestPool 用于复用请求对象 优化GC
type RequestPool[REQ ReqInterface] struct {
	*sync.Pool
}

func NewRequestPool[REQ ReqInterface](newRequest func() REQ) *RequestPool[REQ] {
	r := &RequestPool[REQ]{}
	r.Pool = &sync.Pool{
		New: func() interface{} {
			return newRequest()
		},
	}
	return r
}

func (r *RequestPool[REQ]) NewRequest() REQ {
	return r.Get().(REQ)
}

func (r *RequestPool[REQ]) PutRequest(req REQ) {
	r.Put(req)
}

type Route[REQ ReqInterface, RESP RespInterface] struct {
	RequestPool *RequestPool[REQ]
	Do          func(ctx context.Context, req REQ, opts ...grpc.CallOption) (RESP, error)
}

//var routeMap = map[string]func(ctx context.Context, c *types.UserConn, body IBody) (*pb.ResponseBody, error){}

var httpRouteMap = map[string]gin.HandlerFunc{}
var wsRouteMap = map[string]func(ctx context.Context, connection *gatewayservicelogic.WsConnection, c *pb.GatewayApiRequest) (pb.ResponseCode, []byte, error){}

func AddUnifiedRoute[REQ ReqInterface, RESP RespInterface](svcCtx *svc.ServiceContext, path string, route Route[REQ, RESP]) {
	// http
	AddHttpRoute(path, func(ctx *gin.Context) {
		var response *pb.GatewayApiResponse
		var request REQ
		var err error
		request, response, err = UnifiedHandleHttp(svcCtx, ctx, route)
		defer route.RequestPool.PutRequest(request)
		requestHeader := request.GetHeader()
		if err != nil {
			logx.WithContext(ctx.Request.Context()).Errorf("UnifiedHandleHttp: %s, error: %v", path, err)
			if response == nil {
				response = &pb.GatewayApiResponse{
					Header: i18n.NewServerError(requestHeader, svcCtx.Config.Mode, err),
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
		var request REQ
		requestHeader := connection.Header
		request, response, err = UnifiedHandleWs(svcCtx, ctx, connection, apiRequest, route)
		defer route.RequestPool.PutRequest(request)
		if err != nil {
			logx.WithContext(ctx).Errorf("UnifiedHandleWs: %s, error: %v", path, err)
			if response == nil {
				response = &pb.GatewayApiResponse{
					Header: i18n.NewServerError(requestHeader, svcCtx.Config.Mode, err),
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
			RequestPool: NewRequestPool(func() *pb.GatewayGetUserConnectionReq {
				return &pb.GatewayGetUserConnectionReq{}
			}),
			Do: svcCtx.GatewayService().GatewayGetUserConnection,
		})
		// GatewayBatchGetUserConnectionReq GatewayBatchGetUserConnectionResp
		AddUnifiedRoute(svcCtx, "/v1/gateway/batchGetUserConnection", Route[*pb.GatewayBatchGetUserConnectionReq, *pb.GatewayBatchGetUserConnectionResp]{
			RequestPool: NewRequestPool(func() *pb.GatewayBatchGetUserConnectionReq {
				return &pb.GatewayBatchGetUserConnectionReq{}
			}),
			Do: svcCtx.GatewayService().GatewayBatchGetUserConnection,
		})
		// GatewayGetConnectionByFilterReq GatewayGetConnectionByFilterResp
		AddUnifiedRoute(svcCtx, "/v1/gateway/getConnectionByFilter", Route[*pb.GatewayGetConnectionByFilterReq, *pb.GatewayGetConnectionByFilterResp]{
			RequestPool: NewRequestPool(func() *pb.GatewayGetConnectionByFilterReq {
				return &pb.GatewayGetConnectionByFilterReq{}
			}),
			Do: svcCtx.GatewayService().GatewayGetConnectionByFilter,
		})
		// GatewayWriteDataToWsReq GatewayWriteDataToWsResp
		AddUnifiedRoute(svcCtx, "/v1/gateway/writeDataToWs", Route[*pb.GatewayWriteDataToWsReq, *pb.GatewayWriteDataToWsResp]{
			RequestPool: NewRequestPool(func() *pb.GatewayWriteDataToWsReq {
				return &pb.GatewayWriteDataToWsReq{}
			}),
			Do: svcCtx.GatewayService().GatewayWriteDataToWs,
		})
		// GatewayKickWsReq GatewayKickWsResp
		AddUnifiedRoute(svcCtx, "/v1/gateway/kickWs", Route[*pb.GatewayKickWsReq, *pb.GatewayKickWsResp]{
			RequestPool: NewRequestPool(func() *pb.GatewayKickWsReq {
				return &pb.GatewayKickWsReq{}
			}),
			Do: svcCtx.GatewayService().GatewayKickWs,
		})
		// GatewayKeepAliveReq GatewayKeepAliveResp
		AddWsRoute(svcCtx, "/v1/gateway/white/keepAlive", KeepAliveHandler(svcCtx))
	}
	// user api
	{
		// UserRegisterReq UserRegisterResp
		AddUnifiedRoute(svcCtx, "/v1/user/white/userRegister", Route[*pb.UserRegisterReq, *pb.UserRegisterResp]{
			RequestPool: NewRequestPool(func() *pb.UserRegisterReq {
				return &pb.UserRegisterReq{}
			}),
			Do: svcCtx.AccountService.UserRegister,
		})
		// UserAccessTokenReq UserAccessTokenResp
		AddUnifiedRoute(svcCtx, "/v1/user/white/userAccessToken", Route[*pb.UserAccessTokenReq, *pb.UserAccessTokenResp]{
			RequestPool: NewRequestPool(func() *pb.UserAccessTokenReq {
				return &pb.UserAccessTokenReq{}
			}),
			Do: svcCtx.AccountService.UserAccessToken,
		})
		// CreateRobotReq CreateRobotResp
		AddUnifiedRoute(svcCtx, "/v1/user/createRobot", Route[*pb.CreateRobotReq, *pb.CreateRobotResp]{
			RequestPool: NewRequestPool(func() *pb.CreateRobotReq {
				return &pb.CreateRobotReq{}
			}),
			Do: svcCtx.AccountService.CreateRobot,
		})
		// RefreshUserAccessTokenReq RefreshUserAccessTokenResp
		AddUnifiedRoute(svcCtx, "/v1/user/refreshUserAccessToken", Route[*pb.RefreshUserAccessTokenReq, *pb.RefreshUserAccessTokenResp]{
			RequestPool: NewRequestPool(func() *pb.RefreshUserAccessTokenReq {
				return &pb.RefreshUserAccessTokenReq{}
			}),
			Do: svcCtx.AccountService.RefreshUserAccessToken,
		})
		// RevokeUserAccessTokenReq RevokeUserAccessTokenResp
		AddUnifiedRoute(svcCtx, "/v1/user/revokeUserAccessToken", Route[*pb.RevokeUserAccessTokenReq, *pb.RevokeUserAccessTokenResp]{
			RequestPool: NewRequestPool(func() *pb.RevokeUserAccessTokenReq {
				return &pb.RevokeUserAccessTokenReq{}
			}),
			Do: svcCtx.AccountService.RevokeUserAccessToken,
		})
	}
	// friend api
	{
		//FriendApplyReq FriendApplyResp
		AddUnifiedRoute(svcCtx, "/v1/friend/friendApply", Route[*pb.FriendApplyReq, *pb.FriendApplyResp]{
			RequestPool: NewRequestPool(func() *pb.FriendApplyReq {
				return &pb.FriendApplyReq{}
			}),
			Do: svcCtx.FriendService.FriendApply,
		})
		//ListFriendApplyReq ListFriendApplyResp
		AddUnifiedRoute(svcCtx, "/v1/friend/listFriendApply", Route[*pb.ListFriendApplyReq, *pb.ListFriendApplyResp]{
			RequestPool: NewRequestPool(func() *pb.ListFriendApplyReq {
				return &pb.ListFriendApplyReq{}
			}),
			Do: svcCtx.FriendService.ListFriendApply,
		})
		//FriendApplyHandleReq FriendApplyHandleResp
		AddUnifiedRoute(svcCtx, "/v1/friend/friendApplyHandle", Route[*pb.FriendApplyHandleReq, *pb.FriendApplyHandleResp]{
			RequestPool: NewRequestPool(func() *pb.FriendApplyHandleReq {
				return &pb.FriendApplyHandleReq{}
			}),
			Do: svcCtx.FriendService.FriendApplyHandle,
		})
	}
	// group api
	{
		//GroupCreateReq GroupCreateResp
		AddUnifiedRoute(svcCtx, "/v1/group/groupCreate", Route[*pb.GroupCreateReq, *pb.GroupCreateResp]{
			RequestPool: NewRequestPool(func() *pb.GroupCreateReq {
				return &pb.GroupCreateReq{}
			}),
			Do: svcCtx.GroupService.GroupCreate,
		})
	}
	// notice api
	{
		//ListNoticeReq ListNoticeResp
		AddUnifiedRoute(svcCtx, "/v1/notice/listNotice", Route[*pb.ListNoticeReq, *pb.ListNoticeResp]{
			RequestPool: NewRequestPool(func() *pb.ListNoticeReq {
				return &pb.ListNoticeReq{}
			}),
			Do: svcCtx.NoticeService.ListNotice,
		})
	}
	// message api
	{
		//MessageBatchSendReq MessageBatchSendResp
		AddUnifiedRoute(svcCtx, "/v1/message/messageBatchSend", Route[*pb.MessageBatchSendReq, *pb.MessageBatchSendResp]{
			RequestPool: NewRequestPool(func() *pb.MessageBatchSendReq {
				return &pb.MessageBatchSendReq{}
			}),
			Do: svcCtx.MessageService.MessageBatchSend,
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
