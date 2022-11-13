package handler

import (
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
				{
					Method:  http.MethodGet,
					Path:    "/ping",
					Handler: PingHandler(serverCtx),
				},
			}...,
		),
	)
}
