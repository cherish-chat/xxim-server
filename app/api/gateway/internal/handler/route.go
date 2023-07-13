package handler

import (
	"context"
	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/api/gateway/internal/handler/connectionhandler"
	internalservicelogic "github.com/cherish-chat/xxim-server/app/api/gateway/internal/logic/connectionmanager"
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
	setupGatewayRoutes(svcCtx, engine)
	setupAccountRoutes(svcCtx, engine)
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

*/

func setupAccountRoutes(svcCtx *svc.ServiceContext, engine *gin.Engine) {
	// user api
	{
		// UserRegisterReq UserRegisterResp
		AddUnifiedRoute(svcCtx, "/v1/account/white/userRegister", Route[*peerpb.UserRegisterReq, *peerpb.UserRegisterResp]{
			RequestPool: NewRequestPool(func() *peerpb.UserRegisterReq {
				return &peerpb.UserRegisterReq{}
			}),
			ResponsePool: NewResponsePool(func() *peerpb.UserRegisterResp {
				return &peerpb.UserRegisterResp{}
			}),
			Do: svcCtx.AccountService.UserRegister,
		})
		// UserTokenReq UserTokenResp
		AddUnifiedRoute(svcCtx, "/v1/account/white/userToken", Route[*peerpb.UserTokenReq, *peerpb.UserTokenResp]{
			RequestPool: NewRequestPool(func() *peerpb.UserTokenReq {
				return &peerpb.UserTokenReq{}
			}),
			ResponsePool: NewResponsePool(func() *peerpb.UserTokenResp {
				return &peerpb.UserTokenResp{}
			}),
			Do: svcCtx.AccountService.UserToken,
		})
		// RefreshUserTokenReq RefreshUserTokenResp
		AddUnifiedRoute(svcCtx, "/v1/account/refreshUserToken", Route[*peerpb.RefreshUserTokenReq, *peerpb.RefreshUserTokenResp]{
			RequestPool: NewRequestPool(func() *peerpb.RefreshUserTokenReq {
				return &peerpb.RefreshUserTokenReq{}
			}),
			ResponsePool: NewResponsePool(func() *peerpb.RefreshUserTokenResp {
				return &peerpb.RefreshUserTokenResp{}
			}),
			Do: svcCtx.AccountService.RefreshUserToken,
		})
		// RevokeUserTokenReq RevokeUserTokenResp
		AddUnifiedRoute(svcCtx, "/v1/account/revokeUserToken", Route[*peerpb.RevokeUserTokenReq, *peerpb.RevokeUserTokenResp]{
			RequestPool: NewRequestPool(func() *peerpb.RevokeUserTokenReq {
				return &peerpb.RevokeUserTokenReq{}
			}),
			ResponsePool: NewResponsePool(func() *peerpb.RevokeUserTokenResp {
				return &peerpb.RevokeUserTokenResp{}
			}),
			Do: svcCtx.AccountService.RevokeUserToken,
		})
	}
}

func setupGatewayRoutes(svcCtx *svc.ServiceContext, engine *gin.Engine) {
	// gateway api
	{
		// ListLongConnectionReq ListLongConnectionResp
		AddUnifiedRoute(svcCtx, "/v1/gateway/listLongConnection", Route[*peerpb.ListLongConnectionReq, *peerpb.ListLongConnectionResp]{
			RequestPool: NewRequestPool(func() *peerpb.ListLongConnectionReq {
				return &peerpb.ListLongConnectionReq{}
			}),
			ResponsePool: NewResponsePool(func() *peerpb.ListLongConnectionResp {
				return &peerpb.ListLongConnectionResp{}
			}),
			Do: svcCtx.ConnectionService.ListLongConnection,
		})
		// GatewayKickLongConnectionReq GatewayKickLongConnectionResp
		AddUnifiedRoute(svcCtx, "/v1/gateway/kickLongConnection", Route[*peerpb.GatewayKickLongConnectionReq, *peerpb.GatewayKickLongConnectionResp]{
			RequestPool: NewRequestPool(func() *peerpb.GatewayKickLongConnectionReq {
				return &peerpb.GatewayKickLongConnectionReq{}
			}),
			ResponsePool: NewResponsePool(func() *peerpb.GatewayKickLongConnectionResp {
				return &peerpb.GatewayKickLongConnectionResp{}
			}),
		})
		// GatewayKeepAliveReq GatewayKeepAliveResp
		AddUnifiedRoute(svcCtx, "/v1/gateway/keepAlive", Route[*peerpb.GatewayKeepAliveReq, *peerpb.GatewayKeepAliveResp]{
			RequestPool: NewRequestPool(func() *peerpb.GatewayKeepAliveReq {
				return &peerpb.GatewayKeepAliveReq{}
			}),
			ResponsePool: NewResponsePool(func() *peerpb.GatewayKeepAliveResp {
				return &peerpb.GatewayKeepAliveResp{}
			}),
			Do: svcCtx.InterfaceService.GatewayKeepAlive,
		})

	}
	// 特殊
	{
		connectionHandler := connectionhandler.NewConnectionHandler(svcCtx)
		universalRouteMap["/v1/gateway/white/verifyConnection"] = connectionHandler.VerifyConnection
		universalRouteMap["/v1/gateway/white/authenticationConnection"] = connectionHandler.AuthenticationConnection
	}
}
