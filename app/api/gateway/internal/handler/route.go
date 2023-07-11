package handler

import (
	"context"
	"github.com/cherish-chat/xxim-proto/peerpb"
	internalservicelogic "github.com/cherish-chat/xxim-server/app/api/gateway/internal/logic/internalservice"
	"github.com/cherish-chat/xxim-server/app/api/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/app/api/gateway/internal/types"
	"github.com/gin-gonic/gin"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"sync"
)

type ReqInterface interface {
	proto.Message
	SetHeader(*peerpb.RequestHeader)
	GetHeader() *peerpb.RequestHeader
}

type RespInterface interface {
	proto.Message
	SetHeader(*peerpb.ResponseHeader)
	GetHeader() *peerpb.ResponseHeader
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

var universalRouteMap = map[string]func(ctx context.Context, connection *internalservicelogic.Connection, c *peerpb.GatewayApiRequest) (peerpb.ResponseCode, []byte, error){}

func AddUnifiedRoute[REQ ReqInterface, RESP RespInterface](svcCtx *svc.ServiceContext, path string, route Route[REQ, RESP]) {
	// universal
	universalRouteMap[path] = func(ctx context.Context, connection *internalservicelogic.Connection, apiRequest *peerpb.GatewayApiRequest) (peerpb.ResponseCode, []byte, error) {
		var response *peerpb.GatewayApiResponse
		var err error
		var request REQ
		var respFromPool RESP
		request, respFromPool, response, err = UnifiedHandleUniversal(svcCtx, ctx, connection, apiRequest, route)
		defer route.ResponsePool.PutResponse(respFromPool)
		defer route.RequestPool.PutRequest(request)
		if err != nil {
			logx.WithContext(ctx).Errorf("UnifiedHandleUniversal: %s, error: %v", path, err)
			if response == nil {
				response = &peerpb.GatewayApiResponse{
					Header: peerpb.NewServerError(svcCtx.Config.Mode, err),
					Body:   nil,
				}
			}
			return response.Header.Code, types.MarshalWriteData(response), nil
		}
		if response == nil {
			response = &peerpb.GatewayApiResponse{
				Header: peerpb.NewOkHeader(),
				Body:   nil,
			}
		}
		if response.GetHeader() == nil {
			response.Header = peerpb.NewOkHeader()
		}
		return response.Header.Code, types.MarshalWriteData(response), nil
	}
}

func SetupRoutes(svcCtx *svc.ServiceContext, engine *gin.Engine) {
	//setupGatewayRoutes(svcCtx, engine)
	//setupUserRoutes(svcCtx, engine)
	//setupFriendRoutes(svcCtx, engine)
	//setupGroupRoutes(svcCtx, engine)
	//setupNoticeRoutes(svcCtx, engine)
	//setupMessageRoutes(svcCtx, engine)
}

/*
func setupMessageRoutes(svcCtx *svc.ServiceContext, engine *gin.Engine) {
	// message api
	{
		//MessageBatchSendReq MessageBatchSendResp
		AddUnifiedRoute(svcCtx, "/v1/message/messageBatchSend", Route[*peerpb.MessageBatchSendReq, *peerpb.MessageBatchSendResp]{
			RequestPool: NewRequestPool(func() *peerpb.MessageBatchSendReq {
				return &peerpb.MessageBatchSendReq{}
			}),
			ResponsePool: NewResponsePool(func() *peerpb.MessageBatchSendResp {
				return &peerpb.MessageBatchSendResp{}
			}),
			Do: svcCtx.MessageService.MessageBatchSend,
		})
	}
}

func setupNoticeRoutes(svcCtx *svc.ServiceContext, engine *gin.Engine) {
	// notice api
	{
		//ListNoticeReq ListNoticeResp
		AddUnifiedRoute(svcCtx, "/v1/notice/listNotice", Route[*peerpb.ListNoticeReq, *peerpb.ListNoticeResp]{
			RequestPool: NewRequestPool(func() *peerpb.ListNoticeReq {
				return &peerpb.ListNoticeReq{}
			}),
			ResponsePool: NewResponsePool(func() *peerpb.ListNoticeResp {
				return &peerpb.ListNoticeResp{}
			}),
			Do: svcCtx.NoticeService.ListNotice,
		})
	}
}

func setupGroupRoutes(svcCtx *svc.ServiceContext, engine *gin.Engine) {
	// group api
	{
		//GroupCreateReq GroupCreateResp
		AddUnifiedRoute(svcCtx, "/v1/group/groupCreate", Route[*peerpb.GroupCreateReq, *peerpb.GroupCreateResp]{
			RequestPool: NewRequestPool(func() *peerpb.GroupCreateReq {
				return &peerpb.GroupCreateReq{}
			}),
			ResponsePool: NewResponsePool(func() *peerpb.GroupCreateResp {
				return &peerpb.GroupCreateResp{}
			}),
			Do: svcCtx.GroupService.GroupCreate,
		})
	}
}

func setupFriendRoutes(svcCtx *svc.ServiceContext, engine *gin.Engine) {
	// friend api
	{
		//FriendApplyReq FriendApplyResp
		AddUnifiedRoute(svcCtx, "/v1/friend/friendApply", Route[*peerpb.FriendApplyReq, *peerpb.FriendApplyResp]{
			RequestPool: NewRequestPool(func() *peerpb.FriendApplyReq {
				return &peerpb.FriendApplyReq{}
			}),
			ResponsePool: NewResponsePool(func() *peerpb.FriendApplyResp {
				return &peerpb.FriendApplyResp{}
			}),
			Do: svcCtx.FriendService.FriendApply,
		})
		//ListFriendApplyReq ListFriendApplyResp
		AddUnifiedRoute(svcCtx, "/v1/friend/listFriendApply", Route[*peerpb.ListFriendApplyReq, *peerpb.ListFriendApplyResp]{
			RequestPool: NewRequestPool(func() *peerpb.ListFriendApplyReq {
				return &peerpb.ListFriendApplyReq{}
			}),
			ResponsePool: NewResponsePool(func() *peerpb.ListFriendApplyResp {
				return &peerpb.ListFriendApplyResp{}
			}),
			Do: svcCtx.FriendService.ListFriendApply,
		})
		//FriendApplyHandleReq FriendApplyHandleResp
		AddUnifiedRoute(svcCtx, "/v1/friend/friendApplyHandle", Route[*peerpb.FriendApplyHandleReq, *peerpb.FriendApplyHandleResp]{
			RequestPool: NewRequestPool(func() *peerpb.FriendApplyHandleReq {
				return &peerpb.FriendApplyHandleReq{}
			}),
			ResponsePool: NewResponsePool(func() *peerpb.FriendApplyHandleResp {
				return &peerpb.FriendApplyHandleResp{}
			}),
			Do: svcCtx.FriendService.FriendApplyHandle,
		})
	}
}

func setupUserRoutes(svcCtx *svc.ServiceContext, engine *gin.Engine) {
	// user api
	{
		// UserRegisterReq UserRegisterResp
		AddUnifiedRoute(svcCtx, "/v1/user/white/userRegister", Route[*peerpb.UserRegisterReq, *peerpb.UserRegisterResp]{
			RequestPool: NewRequestPool(func() *peerpb.UserRegisterReq {
				return &peerpb.UserRegisterReq{}
			}),
			ResponsePool: NewResponsePool(func() *peerpb.UserRegisterResp {
				return &peerpb.UserRegisterResp{}
			}),
			Do: svcCtx.AccountService.UserRegister,
		})
		// UserAccessTokenReq UserAccessTokenResp
		AddUnifiedRoute(svcCtx, "/v1/user/white/userAccessToken", Route[*peerpb.UserAccessTokenReq, *peerpb.UserAccessTokenResp]{
			RequestPool: NewRequestPool(func() *peerpb.UserAccessTokenReq {
				return &peerpb.UserAccessTokenReq{}
			}),
			ResponsePool: NewResponsePool(func() *peerpb.UserAccessTokenResp {
				return &peerpb.UserAccessTokenResp{}
			}),
			Do: svcCtx.AccountService.UserAccessToken,
		})
		// CreateRobotReq CreateRobotResp
		AddUnifiedRoute(svcCtx, "/v1/user/createRobot", Route[*peerpb.CreateRobotReq, *peerpb.CreateRobotResp]{
			RequestPool: NewRequestPool(func() *peerpb.CreateRobotReq {
				return &peerpb.CreateRobotReq{}
			}),
			ResponsePool: NewResponsePool(func() *peerpb.CreateRobotResp {
				return &peerpb.CreateRobotResp{}
			}),
			Do: svcCtx.AccountService.CreateRobot,
		})
		// RefreshUserAccessTokenReq RefreshUserAccessTokenResp
		AddUnifiedRoute(svcCtx, "/v1/user/refreshUserAccessToken", Route[*peerpb.RefreshUserAccessTokenReq, *peerpb.RefreshUserAccessTokenResp]{
			RequestPool: NewRequestPool(func() *peerpb.RefreshUserAccessTokenReq {
				return &peerpb.RefreshUserAccessTokenReq{}
			}),
			ResponsePool: NewResponsePool(func() *peerpb.RefreshUserAccessTokenResp {
				return &peerpb.RefreshUserAccessTokenResp{}
			}),
			Do: svcCtx.AccountService.RefreshUserAccessToken,
		})
		// RevokeUserAccessTokenReq RevokeUserAccessTokenResp
		AddUnifiedRoute(svcCtx, "/v1/user/revokeUserAccessToken", Route[*peerpb.RevokeUserAccessTokenReq, *peerpb.RevokeUserAccessTokenResp]{
			RequestPool: NewRequestPool(func() *peerpb.RevokeUserAccessTokenReq {
				return &peerpb.RevokeUserAccessTokenReq{}
			}),
			ResponsePool: NewResponsePool(func() *peerpb.RevokeUserAccessTokenResp {
				return &peerpb.RevokeUserAccessTokenResp{}
			}),
			Do: svcCtx.AccountService.RevokeUserAccessToken,
		})
	}
}

func setupGatewayRoutes(svcCtx *svc.ServiceContext, engine *gin.Engine) {
	// gateway api
	{
		// GatewayGetUserConnectionReq GatewayGetUserConnectionResp
		AddUnifiedRoute(svcCtx, "/v1/gateway/getUserConnection", Route[*peerpb.GatewayGetUserConnectionReq, *peerpb.GatewayGetUserConnectionResp]{
			RequestPool: NewRequestPool(func() *peerpb.GatewayGetUserConnectionReq {
				return &peerpb.GatewayGetUserConnectionReq{}
			}),
			ResponsePool: NewResponsePool(func() *peerpb.GatewayGetUserConnectionResp {
				return &peerpb.GatewayGetUserConnectionResp{}
			}),
			Do: svcCtx.GatewayService().GatewayGetUserConnection,
		})
		// GatewayBatchGetUserConnectionReq GatewayBatchGetUserConnectionResp
		AddUnifiedRoute(svcCtx, "/v1/gateway/batchGetUserConnection", Route[*peerpb.GatewayBatchGetUserConnectionReq, *peerpb.GatewayBatchGetUserConnectionResp]{
			RequestPool: NewRequestPool(func() *peerpb.GatewayBatchGetUserConnectionReq {
				return &peerpb.GatewayBatchGetUserConnectionReq{}
			}),
			ResponsePool: NewResponsePool(func() *peerpb.GatewayBatchGetUserConnectionResp {
				return &peerpb.GatewayBatchGetUserConnectionResp{}
			}),
			Do: svcCtx.GatewayService().GatewayBatchGetUserConnection,
		})
		// GatewayGetConnectionByFilterReq GatewayGetConnectionByFilterResp
		AddUnifiedRoute(svcCtx, "/v1/gateway/getConnectionByFilter", Route[*peerpb.GatewayGetConnectionByFilterReq, *peerpb.GatewayGetConnectionByFilterResp]{
			RequestPool: NewRequestPool(func() *peerpb.GatewayGetConnectionByFilterReq {
				return &peerpb.GatewayGetConnectionByFilterReq{}
			}),
			ResponsePool: NewResponsePool(func() *peerpb.GatewayGetConnectionByFilterResp {
				return &peerpb.GatewayGetConnectionByFilterResp{}
			}),
			Do: svcCtx.GatewayService().GatewayGetConnectionByFilter,
		})
		// GatewayWriteDataToWsReq GatewayWriteDataToWsResp
		AddUnifiedRoute(svcCtx, "/v1/gateway/writeDataToWs", Route[*peerpb.GatewayWriteDataToWsReq, *peerpb.GatewayWriteDataToWsResp]{
			RequestPool: NewRequestPool(func() *peerpb.GatewayWriteDataToWsReq {
				return &peerpb.GatewayWriteDataToWsReq{}
			}),
			ResponsePool: NewResponsePool(func() *peerpb.GatewayWriteDataToWsResp {
				return &peerpb.GatewayWriteDataToWsResp{}
			}),
			Do: svcCtx.GatewayService().GatewayWriteDataToWs,
		})
		// GatewayKickWsReq GatewayKickWsResp
		AddUnifiedRoute(svcCtx, "/v1/gateway/kickWs", Route[*peerpb.GatewayKickWsReq, *peerpb.GatewayKickWsResp]{
			RequestPool: NewRequestPool(func() *peerpb.GatewayKickWsReq {
				return &peerpb.GatewayKickWsReq{}
			}),
			ResponsePool: NewResponsePool(func() *peerpb.GatewayKickWsResp {
				return &peerpb.GatewayKickWsResp{}
			}),
			Do: svcCtx.GatewayService().GatewayKickWs,
		})
		// GatewayKeepAliveReq GatewayKeepAliveResp
		AddUnifiedRoute(svcCtx, "/v1/gateway/keepAlive", Route[*peerpb.GatewayKeepAliveReq, *peerpb.GatewayKeepAliveResp]{
			RequestPool: NewRequestPool(func() *peerpb.GatewayKeepAliveReq {
				return &peerpb.GatewayKeepAliveReq{}
			}),
			ResponsePool: NewResponsePool(func() *peerpb.GatewayKeepAliveResp {
				return &peerpb.GatewayKeepAliveResp{}
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
*/
