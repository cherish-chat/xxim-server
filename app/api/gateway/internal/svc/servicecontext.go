package svc

import (
	"github.com/cherish-chat/xxim-server/app/service/user/client/callbackservice"
	"github.com/cherish-chat/xxim-server/config"
)

type ServiceContext struct {
	Config          config.Config
	CallbackService callbackservice.CallbackService
}

func NewServiceContext(c config.Config) *ServiceContext {
	userClient, _ := c.GetUserRpcClient()
	return &ServiceContext{
		Config:          c,
		CallbackService: callbackservice.NewCallbackService(userClient),
	}
}
