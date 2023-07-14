package svc

import (
	"github.com/cherish-chat/xxim-server/app/api/gateway/client/connectionservice"
	"github.com/cherish-chat/xxim-server/app/api/gateway/client/interfaceservice"
	"github.com/cherish-chat/xxim-server/app/service/conversation/client/friendservice"
	"github.com/cherish-chat/xxim-server/app/service/message/client/messageservice"
	"github.com/cherish-chat/xxim-server/app/service/user/client/accountservice"
	"github.com/cherish-chat/xxim-server/app/service/user/client/callbackservice"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/config"
)

type ServiceContext struct {
	Config          config.Config
	CallbackService callbackservice.CallbackService
	AccountService  accountservice.AccountService

	FriendService friendservice.FriendService

	ConnectionService connectionservice.ConnectionService
	InterfaceService  interfaceservice.InterfaceService

	MessageService messageservice.MessageService
	RsaInstance    *utils.XRsa
}

func NewServiceContext(c config.Config) *ServiceContext {
	userClient, _ := c.GetUserRpcClient()
	conversationClient, _ := c.GetConversationRpcClient()
	gatewayClient, _ := c.GetGatewayRpcClient()
	messageClient, _ := c.GetMessageRpcClient()
	return &ServiceContext{
		Config:          c,
		CallbackService: callbackservice.NewCallbackService(userClient),
		AccountService:  accountservice.NewAccountService(userClient),

		FriendService: friendservice.NewFriendService(conversationClient),

		ConnectionService: connectionservice.NewConnectionService(gatewayClient),
		InterfaceService:  interfaceservice.NewInterfaceService(gatewayClient),

		MessageService: messageservice.NewMessageService(messageClient),

		RsaInstance: utils.NewRsa(c.GetPublicKey(), c.GetPrivateKey()),
	}
}
