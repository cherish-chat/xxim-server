package svc

import (
	"github.com/cherish-chat/xxim-server/app/conn/internal/config"
	"github.com/cherish-chat/xxim-server/app/im/imservice"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config    config.Config
	Redis     *redis.Redis
	Mgo       *xmgo.Client
	imService imservice.ImService
	PodIp     string
}

func NewServiceContext(c config.Config) *ServiceContext {
	s := &ServiceContext{
		Config: c,
		PodIp:  utils.GetPodIp(),
	}
	s.mustSetup()
	return s
}

func (s *ServiceContext) mustSetup() {
	s.Redis = s.Config.Redis.NewRedis()
	s.Mgo = xmgo.NewClient(s.Config.Mongo)
}

func (s *ServiceContext) ImService() imservice.ImService {
	if s.imService == nil {
		s.imService = imservice.NewImService(zrpc.MustNewClient(s.Config.ImRpc))
	}
	return s.imService
}
