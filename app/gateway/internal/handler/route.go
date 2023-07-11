package handler

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/handler/connectionhandler"
	gatewayservicelogic "github.com/cherish-chat/xxim-server/app/gateway/internal/logic/gatewayservice"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/types"
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
	SetHeader(*pb.ResponseHeader)
	GetHeader() *pb.ResponseHeader
}

// RequestPool 用于复用请求对象 优化GC
type RequestPool[REQ ReqInterface] struct {
	*sync.Pool
}

// ResponsePool 用于复用响应对象 优化GC
type ResponsePool[RESP RespInterface] struct {
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

func NewResponsePool[RESP RespInterface](newResponse func() RESP) *ResponsePool[RESP] {
	r := &ResponsePool[RESP]{}
	r.Pool = &sync.Pool{
		New: func() interface{} {
			return newResponse()
		},
	}
	return r
}

func (r *RequestPool[REQ]) NewRequest() REQ {
	return r.Get().(REQ)
}

func (r *ResponsePool[RESP]) NewResponse() RESP {
	return r.Get().(RESP)
}

func (r *RequestPool[REQ]) PutRequest(req REQ) {
	r.Put(req)
}

func (r *ResponsePool[RESP]) PutResponse(resp RESP) {
	r.Put(resp)
}

type Route[REQ ReqInterface, RESP RespInterface] struct {
	RequestPool  *RequestPool[REQ]
	ResponsePool *ResponsePool[RESP]
	Do           func(ctx context.Context, req REQ, opts ...grpc.CallOption) (RESP, error)
}

var universalRouteMap = map[string]func(ctx context.Context, connection *gatewayservicelogic.Connection, c *pb.GatewayApiRequest) (pb.ResponseCode, []byte, error){}

func AddUnifiedRoute[REQ ReqInterface, RESP RespInterface](svcCtx *svc.ServiceContext, path string, route Route[REQ, RESP]) {
	// universal
	universalRouteMap[path] = func(ctx context.Context, connection *gatewayservicelogic.Connection, apiRequest *pb.GatewayApiRequest) (pb.ResponseCode, []byte, error) {
		var response *pb.GatewayApiResponse
		var err error
		var request REQ
		var respFromPool RESP
		request, respFromPool, response, err = UnifiedHandleUniversal(svcCtx, ctx, connection, apiRequest, route)
		defer route.ResponsePool.PutResponse(respFromPool)
		defer route.RequestPool.PutRequest(request)
		if err != nil {
			logx.WithContext(ctx).Errorf("UnifiedHandleUniversal: %s, error: %v", path, err)
			if response == nil {
				response = &pb.GatewayApiResponse{
					Header: i18n.NewServerError(svcCtx.Config.Mode, err),
					Body:   nil,
				}
			}
			return response.Header.Code, types.MarshalWriteData(response), nil
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
		return response.Header.Code, types.MarshalWriteData(response), nil
	}
}

func SetupRoutes(svcCtx *svc.ServiceContext, engine *gin.Engine) {
	setupGatewayRoutes(svcCtx, engine)
	setupUserRoutes(svcCtx, engine)
	setupFriendRoutes(svcCtx, engine)
	setupGroupRoutes(svcCtx, engine)
	setupNoticeRoutes(svcCtx, engine)
	setupMessageRoutes(svcCtx, engine)

	wsHandler := NewWsHandler(svcCtx)
	engine.GET("/ws", wsHandler.Upgrade)
}

func setupMessageRoutes(svcCtx *svc.ServiceContext, engine *gin.Engine) {
	// message api
	{
		//MessageBatchSendReq MessageBatchSendResp
		AddUnifiedRoute(svcCtx, "/v1/message/messageBatchSend", Route[*pb.MessageBatchSendReq, *pb.MessageBatchSendResp]{
			RequestPool: NewRequestPool(func() *pb.MessageBatchSendReq {
				return &pb.MessageBatchSendReq{}
			}),
			ResponsePool: NewResponsePool(func() *pb.MessageBatchSendResp {
				return &pb.MessageBatchSendResp{}
			}),
			Do: svcCtx.MessageService.MessageBatchSend,
		})
	}
}

func setupNoticeRoutes(svcCtx *svc.ServiceContext, engine *gin.Engine) {
	// notice api
	{
		//ListNoticeReq ListNoticeResp
		AddUnifiedRoute(svcCtx, "/v1/notice/listNotice", Route[*pb.ListNoticeReq, *pb.ListNoticeResp]{
			RequestPool: NewRequestPool(func() *pb.ListNoticeReq {
				return &pb.ListNoticeReq{}
			}),
			ResponsePool: NewResponsePool(func() *pb.ListNoticeResp {
				return &pb.ListNoticeResp{}
			}),
			Do: svcCtx.NoticeService.ListNotice,
		})
	}
}

func setupGroupRoutes(svcCtx *svc.ServiceContext, engine *gin.Engine) {
	// group api
	{
		//GroupCreateReq GroupCreateResp
		AddUnifiedRoute(svcCtx, "/v1/group/groupCreate", Route[*pb.GroupCreateReq, *pb.GroupCreateResp]{
			RequestPool: NewRequestPool(func() *pb.GroupCreateReq {
				return &pb.GroupCreateReq{}
			}),
			ResponsePool: NewResponsePool(func() *pb.GroupCreateResp {
				return &pb.GroupCreateResp{}
			}),
			Do: svcCtx.GroupService.GroupCreate,
		})
	}
}

func setupFriendRoutes(svcCtx *svc.ServiceContext, engine *gin.Engine) {
	// friend api
	{
		//FriendApplyReq FriendApplyResp
		AddUnifiedRoute(svcCtx, "/v1/friend/friendApply", Route[*pb.FriendApplyReq, *pb.FriendApplyResp]{
			RequestPool: NewRequestPool(func() *pb.FriendApplyReq {
				return &pb.FriendApplyReq{}
			}),
			ResponsePool: NewResponsePool(func() *pb.FriendApplyResp {
				return &pb.FriendApplyResp{}
			}),
			Do: svcCtx.FriendService.FriendApply,
		})
		//ListFriendApplyReq ListFriendApplyResp
		AddUnifiedRoute(svcCtx, "/v1/friend/listFriendApply", Route[*pb.ListFriendApplyReq, *pb.ListFriendApplyResp]{
			RequestPool: NewRequestPool(func() *pb.ListFriendApplyReq {
				return &pb.ListFriendApplyReq{}
			}),
			ResponsePool: NewResponsePool(func() *pb.ListFriendApplyResp {
				return &pb.ListFriendApplyResp{}
			}),
			Do: svcCtx.FriendService.ListFriendApply,
		})
		//FriendApplyHandleReq FriendApplyHandleResp
		AddUnifiedRoute(svcCtx, "/v1/friend/friendApplyHandle", Route[*pb.FriendApplyHandleReq, *pb.FriendApplyHandleResp]{
			RequestPool: NewRequestPool(func() *pb.FriendApplyHandleReq {
				return &pb.FriendApplyHandleReq{}
			}),
			ResponsePool: NewResponsePool(func() *pb.FriendApplyHandleResp {
				return &pb.FriendApplyHandleResp{}
			}),
			Do: svcCtx.FriendService.FriendApplyHandle,
		})
	}
}

func setupUserRoutes(svcCtx *svc.ServiceContext, engine *gin.Engine) {
	// user api
	{
		// UserRegisterReq UserRegisterResp
		AddUnifiedRoute(svcCtx, "/v1/user/white/userRegister", Route[*pb.UserRegisterReq, *pb.UserRegisterResp]{
			RequestPool: NewRequestPool(func() *pb.UserRegisterReq {
				return &pb.UserRegisterReq{}
			}),
			ResponsePool: NewResponsePool(func() *pb.UserRegisterResp {
				return &pb.UserRegisterResp{}
			}),
			Do: svcCtx.AccountService.UserRegister,
		})
		// UserAccessTokenReq UserAccessTokenResp
		AddUnifiedRoute(svcCtx, "/v1/user/white/userAccessToken", Route[*pb.UserAccessTokenReq, *pb.UserAccessTokenResp]{
			RequestPool: NewRequestPool(func() *pb.UserAccessTokenReq {
				return &pb.UserAccessTokenReq{}
			}),
			ResponsePool: NewResponsePool(func() *pb.UserAccessTokenResp {
				return &pb.UserAccessTokenResp{}
			}),
			Do: svcCtx.AccountService.UserAccessToken,
		})
		// CreateRobotReq CreateRobotResp
		AddUnifiedRoute(svcCtx, "/v1/user/createRobot", Route[*pb.CreateRobotReq, *pb.CreateRobotResp]{
			RequestPool: NewRequestPool(func() *pb.CreateRobotReq {
				return &pb.CreateRobotReq{}
			}),
			ResponsePool: NewResponsePool(func() *pb.CreateRobotResp {
				return &pb.CreateRobotResp{}
			}),
			Do: svcCtx.AccountService.CreateRobot,
		})
		// RefreshUserAccessTokenReq RefreshUserAccessTokenResp
		AddUnifiedRoute(svcCtx, "/v1/user/refreshUserAccessToken", Route[*pb.RefreshUserAccessTokenReq, *pb.RefreshUserAccessTokenResp]{
			RequestPool: NewRequestPool(func() *pb.RefreshUserAccessTokenReq {
				return &pb.RefreshUserAccessTokenReq{}
			}),
			ResponsePool: NewResponsePool(func() *pb.RefreshUserAccessTokenResp {
				return &pb.RefreshUserAccessTokenResp{}
			}),
			Do: svcCtx.AccountService.RefreshUserAccessToken,
		})
		// RevokeUserAccessTokenReq RevokeUserAccessTokenResp
		AddUnifiedRoute(svcCtx, "/v1/user/revokeUserAccessToken", Route[*pb.RevokeUserAccessTokenReq, *pb.RevokeUserAccessTokenResp]{
			RequestPool: NewRequestPool(func() *pb.RevokeUserAccessTokenReq {
				return &pb.RevokeUserAccessTokenReq{}
			}),
			ResponsePool: NewResponsePool(func() *pb.RevokeUserAccessTokenResp {
				return &pb.RevokeUserAccessTokenResp{}
			}),
			Do: svcCtx.AccountService.RevokeUserAccessToken,
		})
	}
}

func setupGatewayRoutes(svcCtx *svc.ServiceContext, engine *gin.Engine) {
	// gateway api
	{
		// GatewayGetUserConnectionReq GatewayGetUserConnectionResp
		AddUnifiedRoute(svcCtx, "/v1/gateway/getUserConnection", Route[*pb.GatewayGetUserConnectionReq, *pb.GatewayGetUserConnectionResp]{
			RequestPool: NewRequestPool(func() *pb.GatewayGetUserConnectionReq {
				return &pb.GatewayGetUserConnectionReq{}
			}),
			ResponsePool: NewResponsePool(func() *pb.GatewayGetUserConnectionResp {
				return &pb.GatewayGetUserConnectionResp{}
			}),
			Do: svcCtx.GatewayService().GatewayGetUserConnection,
		})
		// GatewayBatchGetUserConnectionReq GatewayBatchGetUserConnectionResp
		AddUnifiedRoute(svcCtx, "/v1/gateway/batchGetUserConnection", Route[*pb.GatewayBatchGetUserConnectionReq, *pb.GatewayBatchGetUserConnectionResp]{
			RequestPool: NewRequestPool(func() *pb.GatewayBatchGetUserConnectionReq {
				return &pb.GatewayBatchGetUserConnectionReq{}
			}),
			ResponsePool: NewResponsePool(func() *pb.GatewayBatchGetUserConnectionResp {
				return &pb.GatewayBatchGetUserConnectionResp{}
			}),
			Do: svcCtx.GatewayService().GatewayBatchGetUserConnection,
		})
		// GatewayGetConnectionByFilterReq GatewayGetConnectionByFilterResp
		AddUnifiedRoute(svcCtx, "/v1/gateway/getConnectionByFilter", Route[*pb.GatewayGetConnectionByFilterReq, *pb.GatewayGetConnectionByFilterResp]{
			RequestPool: NewRequestPool(func() *pb.GatewayGetConnectionByFilterReq {
				return &pb.GatewayGetConnectionByFilterReq{}
			}),
			ResponsePool: NewResponsePool(func() *pb.GatewayGetConnectionByFilterResp {
				return &pb.GatewayGetConnectionByFilterResp{}
			}),
			Do: svcCtx.GatewayService().GatewayGetConnectionByFilter,
		})
		// GatewayWriteDataToWsReq GatewayWriteDataToWsResp
		AddUnifiedRoute(svcCtx, "/v1/gateway/writeDataToWs", Route[*pb.GatewayWriteDataToWsReq, *pb.GatewayWriteDataToWsResp]{
			RequestPool: NewRequestPool(func() *pb.GatewayWriteDataToWsReq {
				return &pb.GatewayWriteDataToWsReq{}
			}),
			ResponsePool: NewResponsePool(func() *pb.GatewayWriteDataToWsResp {
				return &pb.GatewayWriteDataToWsResp{}
			}),
			Do: svcCtx.GatewayService().GatewayWriteDataToWs,
		})
		// GatewayKickWsReq GatewayKickWsResp
		AddUnifiedRoute(svcCtx, "/v1/gateway/kickWs", Route[*pb.GatewayKickWsReq, *pb.GatewayKickWsResp]{
			RequestPool: NewRequestPool(func() *pb.GatewayKickWsReq {
				return &pb.GatewayKickWsReq{}
			}),
			ResponsePool: NewResponsePool(func() *pb.GatewayKickWsResp {
				return &pb.GatewayKickWsResp{}
			}),
			Do: svcCtx.GatewayService().GatewayKickWs,
		})
		// GatewayKeepAliveReq GatewayKeepAliveResp
		AddUnifiedRoute(svcCtx, "/v1/gateway/keepAlive", Route[*pb.GatewayKeepAliveReq, *pb.GatewayKeepAliveResp]{
			RequestPool: NewRequestPool(func() *pb.GatewayKeepAliveReq {
				return &pb.GatewayKeepAliveReq{}
			}),
			ResponsePool: NewResponsePool(func() *pb.GatewayKeepAliveResp {
				return &pb.GatewayKeepAliveResp{}
			}),
			Do: svcCtx.GatewayService().GatewayKeepAlive,
		})

	}
	// 特殊
	{
		connectionHandler := connectionhandler.NewConnectionHandler(svcCtx)
		universalRouteMap["/v1/gateway/white/verifyConnection"] = connectionHandler.VerifyConnection
		universalRouteMap["/v1/gateway/white/authenticationConnection"] = connectionHandler.AuthenticationConnection
	}
}
