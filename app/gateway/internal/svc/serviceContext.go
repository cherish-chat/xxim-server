package svc

import (
	"github.com/cherish-chat/xxim-server/app/gateway/gatewayservice"
	"github.com/cherish-chat/xxim-server/app/gateway/internal/config"
	"github.com/cherish-chat/xxim-server/app/user/userservice"
	"github.com/cherish-chat/xxim-server/common/xcache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"strings"
	"time"
)

type ServiceContext struct {
	Config         config.Config
	gatewayService gatewayservice.GatewayService
	UserService    userservice.UserService
	Redis          *redis.Redis
}

func NewServiceContext(c config.Config) *ServiceContext {
	s := &ServiceContext{
		Config: c,
		UserService: userservice.NewUserService(zrpc.MustNewClient(
			c.RpcClientConf.User,
			zrpc.WithNonBlock(),
			zrpc.WithTimeout(time.Duration(c.Timeout)*time.Millisecond),
		)),
		Redis: xcache.MustNewRedis(c.RedisConf),
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
