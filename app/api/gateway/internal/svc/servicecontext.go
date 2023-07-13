package svc

import (
	"github.com/cherish-chat/xxim-server/app/api/gateway/client/connectionservice"
	"github.com/cherish-chat/xxim-server/app/api/gateway/client/interfaceservice"
	"github.com/cherish-chat/xxim-server/app/service/user/client/accountservice"
	"github.com/cherish-chat/xxim-server/app/service/user/client/callbackservice"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/config"
)

type ServiceContext struct {
	Config          config.Config
	CallbackService callbackservice.CallbackService
	AccountService  accountservice.AccountService

	ConnectionService connectionservice.ConnectionService
	InterfaceService  interfaceservice.InterfaceService

	RsaInstance *utils.XRsa
}

func NewServiceContext(c config.Config) *ServiceContext {
	userClient, _ := c.GetUserRpcClient()
	gatewayClient, _ := c.GetGatewayRpcClient()
	return &ServiceContext{
		Config:          c,
		CallbackService: callbackservice.NewCallbackService(userClient),
		AccountService:  accountservice.NewAccountService(userClient),

		ConnectionService: connectionservice.NewConnectionService(gatewayClient),
		InterfaceService:  interfaceservice.NewInterfaceService(gatewayClient),

		RsaInstance: utils.NewRsa(c.GetPublicKey(), c.GetPrivateKey()),
	}
}
