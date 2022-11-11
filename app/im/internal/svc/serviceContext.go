package svc

import (
	"github.com/cherish-chat/xxim-server/app/conn/connservice"
	"github.com/cherish-chat/xxim-server/app/im/internal/config"
	msgservice "github.com/cherish-chat/xxim-server/app/msg/msgService"
	"github.com/cherish-chat/xxim-server/common/utils/ip2region"
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config      config.Config
	ConnPodsMgr *connservice.ConnPodsMgr
	zedis       *redis.Redis
	mongo       *xmgo.Client
	msgService  msgservice.MsgService
}

func (s *ServiceContext) Redis() *redis.Redis {
	if s.zedis == nil {
		s.zedis = s.Config.Redis.NewRedis()
	}
	return s.zedis
}

func (s *ServiceContext) Mongo() *xmgo.Client {
	if s.mongo == nil {
		s.mongo = xmgo.NewClient(s.Config.Mongo)
	}
	return s.mongo
}

func (s *ServiceContext) MsgService() msgservice.MsgService {
	if s.msgService == nil {
		s.msgService = msgservice.NewMsgService(zrpc.MustNewClient(s.Config.MsgRpc))
	}
	return s.msgService
}

func NewServiceContext(c config.Config) *ServiceContext {
	s := &ServiceContext{
		Config: c,
	}
	ip2region.Init(c.Ip2RegionUrl)
	s.ConnPodsMgr = connservice.NewConnPodsMgr(c.ConnRpc)
	return s
}
