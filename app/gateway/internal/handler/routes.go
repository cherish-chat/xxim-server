package handler

import (
	"github.com/cherish-chat/xxim-server/app/gateway/internal/handler/grouphandler"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/handler/msghandler"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/handler/relationhandler"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/handler/userhandler"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/svc"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/wrapper"
	"net/http"

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
				// setSingleChatSetting
				{
					Method:  http.MethodPost,
					Path:    "/setSingleChatSetting",
					Handler: wrapper.WrapHandler(serverCtx, relationhandler.SetSingleChatSettingConfig(serverCtx)),
				},
				// setSingleMsgNotifyOpt
				{
					Method:  http.MethodPost,
					Path:    "/setSingleMsgNotifyOpt",
					Handler: wrapper.WrapHandler(serverCtx, relationhandler.SetSingleMsgNotifyOptConfig(serverCtx)),
				},
				// getSingleChatSetting
				{
					Method:  http.MethodPost,
					Path:    "/getSingleChatSetting",
					Handler: wrapper.WrapHandler(serverCtx, relationhandler.GetSingleChatSettingConfig(serverCtx)),
				},
				// getSingleMsgNotifyOpt
				{
					Method:  http.MethodPost,
					Path:    "/getSingleMsgNotifyOpt",
					Handler: wrapper.WrapHandler(serverCtx, relationhandler.GetSingleMsgNotifyOptConfig(serverCtx)),
				},
				// getFriendList
				{
					Method:  http.MethodPost,
					Path:    "/getFriendList",
					Handler: wrapper.WrapHandler(serverCtx, relationhandler.GetFriendListConfig(serverCtx)),
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
					Path:    "/getMsgListByConvId",
					Handler: wrapper.WrapHandler(serverCtx, msghandler.GetMsgListByConvIdConfig(serverCtx)),
				},
				// getMsgById
				{
					Method:  http.MethodPost,
					Path:    "/getMsgById",
					Handler: wrapper.WrapHandler(serverCtx, msghandler.GetMsgByIdConfig(serverCtx)),
				},
			}...,
		),
		rest.WithPrefix("/v1/msg"),
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
}
