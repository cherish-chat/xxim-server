package svc

import (
	"github.com/cherish-chat/xxim-server/app/conversation/client/friendservice"
	"github.com/cherish-chat/xxim-server/app/gateway/client/gatewayservice"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/config"
	"github.com/cherish-chat/xxim-server/app/message/client/messageservice"
	"github.com/cherish-chat/xxim-server/app/message/client/noticeservice"
	"github.com/cherish-chat/xxim-server/app/user/client/accountservice"
	"github.com/cherish-chat/xxim-server/app/user/client/callbackservice"
	"github.com/cherish-chat/xxim-server/app/user/client/infoservice"
	"github.com/cherish-chat/xxim-server/common/xcache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"strings"
	"time"
)

type ServiceContext struct {
	Config config.Config
	Redis  *redis.Redis

	gatewayService gatewayservice.GatewayService
	//User
	CallbackService callbackservice.CallbackService
	AccountService  accountservice.AccountService
	InfoService     infoservice.InfoService
	//Conversation
	FriendService friendservice.FriendService
	//Message
	NoticeService  noticeservice.NoticeService
	MessageService messageservice.MessageService
}

func NewServiceContext(c config.Config) *ServiceContext {

	userClient := zrpc.MustNewClient(
		c.RpcClientConf.User,
		zrpc.WithNonBlock(),
		zrpc.WithTimeout(time.Duration(c.Timeout)*time.Millisecond),
	)
	conversationClient := zrpc.MustNewClient(
		c.RpcClientConf.Conversation,
		zrpc.WithNonBlock(),
		zrpc.WithTimeout(time.Duration(c.Timeout)*time.Millisecond),
	)
	messageClient := zrpc.MustNewClient(
		c.RpcClientConf.Message,
		zrpc.WithNonBlock(),
		zrpc.WithTimeout(time.Duration(c.Timeout)*time.Millisecond),
	)

	s := &ServiceContext{
		Config:          c,
		CallbackService: callbackservice.NewCallbackService(userClient),
		AccountService:  accountservice.NewAccountService(userClient),
		InfoService:     infoservice.NewInfoService(userClient),
		FriendService:   friendservice.NewFriendService(conversationClient),
		NoticeService:   noticeservice.NewNoticeService(messageClient),
		MessageService:  messageservice.NewMessageService(messageClient),
		Redis:           xcache.MustNewRedis(c.RedisConf),
	}
	return s
}

func (s *ServiceContext) GatewayService() gatewayservice.GatewayService {
	if s.gatewayService == nil {
		listenOnSplit := strings.Split(s.Config.ListenOn, ":")
		rpcPort := listenOnSplit[len(listenOnSplit)-1]
		s.gatewayService = gatewayservice.MustNewGatewayService(zrpc.RpcClientConf{
			Endpoints: []string{"127.0.0.1:" + rpcPort},
			NonBlock:  true,
		})
	}
	return s.gatewayService
}
