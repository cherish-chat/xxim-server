package handler

import (
	"github.com/cherish-chat/xxim-server/app/gateway/internal/handler/grouphandler"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/handler/imhandler"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/handler/msghandler"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/handler/relationhandler"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/handler/userhandler"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/wrapper"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{},
			[]rest.Route{
				// login
				{
					Method:  http.MethodPost,
					Path:    "/white/login",
					Handler: wrapper.WrapHandler(serverCtx, userhandler.LoginConfig(serverCtx)),
				},
				// confirmRegister
				{
					Method:  http.MethodPost,
					Path:    "/white/confirmRegister",
					Handler: wrapper.WrapHandler(serverCtx, userhandler.ConfirmRegisterConfig(serverCtx)),
				},
				// searchUsersByKeyword
				{
					Method:  http.MethodPost,
					Path:    "/searchUsersByKeyword",
					Handler: wrapper.WrapHandler(serverCtx, userhandler.SearchUsersByKeywordConfig(serverCtx)),
				},
				// getUserHome
				{
					Method:  http.MethodPost,
					Path:    "/getUserHome",
					Handler: wrapper.WrapHandler(serverCtx, userhandler.GetUserHomeConfig(serverCtx)),
				},
				// getUserSettings
				{
					Method:  http.MethodPost,
					Path:    "/getUserSettings",
					Handler: wrapper.WrapHandler(serverCtx, userhandler.GetUserSettingsConfig(serverCtx)),
				},
				// setUserSettings
				{
					Method:  http.MethodPost,
					Path:    "/setUserSettings",
					Handler: wrapper.WrapHandler(serverCtx, userhandler.SetUserSettingsConfig(serverCtx)),
				},
			}...,
		),
		rest.WithPrefix("/v1/user"),
	)
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{},
			[]rest.Route{
				// requestAddFriend
				{
					Method:  http.MethodPost,
					Path:    "/requestAddFriend",
					Handler: wrapper.WrapHandler(serverCtx, relationhandler.RequestAddFriendConfig(serverCtx)),
				},
				// acceptAddFriend
				{
					Method:  http.MethodPost,
					Path:    "/acceptAddFriend",
					Handler: wrapper.WrapHandler(serverCtx, relationhandler.AcceptAddFriendConfig(serverCtx)),
				},
				// rejectAddFriend
				{
					Method:  http.MethodPost,
					Path:    "/rejectAddFriend",
					Handler: wrapper.WrapHandler(serverCtx, relationhandler.RejectAddFriendConfig(serverCtx)),
				},
				// blockUser
				{
					Method:  http.MethodPost,
					Path:    "/blockUser",
					Handler: wrapper.WrapHandler(serverCtx, relationhandler.BlockUserConfig(serverCtx)),
				},
				// deleteBlockUser
				{
					Method:  http.MethodPost,
					Path:    "/deleteBlockUser",
					Handler: wrapper.WrapHandler(serverCtx, relationhandler.DeleteBlockUserConfig(serverCtx)),
				},
				// deleteFriend
				{
					Method:  http.MethodPost,
					Path:    "/deleteFriend",
					Handler: wrapper.WrapHandler(serverCtx, relationhandler.DeleteFriendConfig(serverCtx)),
				},
				// setSingleConvSetting
				{
					Method:  http.MethodPost,
					Path:    "/setSingleConvSetting",
					Handler: wrapper.WrapHandler(serverCtx, relationhandler.SetSingleConvSettingConfig(serverCtx)),
				},
				// getSingleConvSetting
				{
					Method:  http.MethodPost,
					Path:    "/getSingleConvSetting",
					Handler: wrapper.WrapHandler(serverCtx, relationhandler.GetSingleConvSettingConfig(serverCtx)),
				},
				// getFriendList
				{
					Method:  http.MethodPost,
					Path:    "/getFriendList",
					Handler: wrapper.WrapHandler(serverCtx, relationhandler.GetFriendListConfig(serverCtx)),
				},
				// getMyFriendEventList
				{
					Method:  http.MethodPost,
					Path:    "/getMyFriendEventList",
					Handler: wrapper.WrapHandler(serverCtx, relationhandler.GetMyFriendEventListConfig(serverCtx)),
				},
			}...,
		),
		rest.WithPrefix("/v1/relation"),
	)
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{},
			[]rest.Route{
				// createGroup
				{
					Method:  http.MethodPost,
					Path:    "/createGroup",
					Handler: wrapper.WrapHandler(serverCtx, grouphandler.CreateGroupConfig(serverCtx)),
				},
			}...,
		),
		rest.WithPrefix("/v1/group"),
	)
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{},
			[]rest.Route{
				// sendMsg
				{
					Method:  http.MethodPost,
					Path:    "/sendMsg",
					Handler: wrapper.WrapHandler(serverCtx, msghandler.SendMsgConfig(serverCtx)),
				},
				// getMsgListByConvId
				{
					Method:  http.MethodPost,
					Path:    "/batchGetMsgListByConvId",
					Handler: wrapper.WrapHandler(serverCtx, msghandler.BatchGetMsgListByConvIdConfig(serverCtx)),
				},
				// getMsgById
				{
					Method:  http.MethodPost,
					Path:    "/getMsgById",
					Handler: wrapper.WrapHandler(serverCtx, msghandler.GetMsgByIdConfig(serverCtx)),
				},
				// batchGetConvSeq
				{
					Method:  http.MethodPost,
					Path:    "/batchGetConvSeq",
					Handler: wrapper.WrapHandler(serverCtx, msghandler.BatchGetConvSeqConfig(serverCtx)),
				},
			}...,
		),
		rest.WithPrefix("/v1/msg"),
	)
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{},
			[]rest.Route{
				// getAppSystemConfig
				{
					Method:  http.MethodPost,
					Path:    "/white/getAppSystemConfig",
					Handler: wrapper.WrapHandler(serverCtx, imhandler.GetAppSystemConfigConfig(serverCtx)),
				},
			}...,
		),
		rest.WithPrefix("/v1/im"),
	)
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/ping",
					Handler: PingHandler(serverCtx),
				},
			}...,
		),
	)
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{},
			[]rest.Route{
				{
					Method:  http.MethodGet,
					Path:    "/ws",
					Handler: WsHandler(serverCtx),
				},
			}...,
		),
	)
	server.AddRoutes(
		rest.WithMiddlewares(
			[]rest.Middleware{},
			[]rest.Route{
				{
					Method:  http.MethodPost,
					Path:    "/authVerify",
					Handler: AuthHandler(serverCtx),
				},
			}...,
		),
	)
	go func() {
		time.Sleep(1 * time.Second)
		server.PrintRoutes()
	}()
}
