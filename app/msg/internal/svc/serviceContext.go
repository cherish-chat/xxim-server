package svc

import (
	"github.com/cherish-chat/xxim-server/app/group/groupservice"
	"github.com/cherish-chat/xxim-server/app/im/imservice"
	"github.com/cherish-chat/xxim-server/app/msg/internal/config"
	"github.com/cherish-chat/xxim-server/app/msg/msgmodel"
	"github.com/cherish-chat/xxim-server/app/notice/noticeservice"
	"github.com/cherish-chat/xxim-server/app/relation/relationservice"
	"github.com/cherish-chat/xxim-server/app/user/userservice"
	"github.com/cherish-chat/xxim-server/common/i18n"
	"github.com/cherish-chat/xxim-server/common/pkg/mobpush"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xconf"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/cherish-chat/xxim-server/common/xtdmq"
	"github.com/zeromicro/go-zero/core/limit"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config             config.Config
	msgProducer        *xtdmq.TDMQProducer
	zedis              *redis.Redis
	zedisSub           *redis.Redis
	mysql              *gorm.DB
	imService          imservice.ImService
	relationService    relationservice.RelationService
	groupService       groupservice.GroupService
	userService        userservice.UserService
	noticeService      noticeservice.NoticeService
	MobPush            *mobpush.Pusher
	ConfigMgr          *xconf.ConfigMgr
	SyncSendMsgLimiter *limit.TokenLimiter
	*i18n.I18N
}

func NewServiceContext(c config.Config) *ServiceContext {
	s := &ServiceContext{
		Config: c,
	}
	s.I18N = i18n.NewI18N(s.Mysql())
	s.Mysql().AutoMigrate(
		msgmodel.Msg{},
		xorm.HashKv{},
	)
	s.MobPush = mobpush.NewPusher(c.MobPush)
	s.ConfigMgr = xconf.NewConfigMgr(s.Mysql(), s.Redis(), "system")
	s.SyncSendMsgLimiter = limit.NewTokenLimiter(c.SyncSendMsgLimit.Rate, c.SyncSendMsgLimit.Burst, s.Redis(), rediskey.SyncSendMsgLimiter())
	return s
}

func (s *ServiceContext) MsgProducer() *xtdmq.TDMQProducer {
	if s.msgProducer == nil {
		s.msgProducer = xtdmq.NewTDMQProducer(s.Config.TDMQ.TDMQConfig, s.Config.TDMQ.Producer)
	}
	return s.msgProducer
}

func (s *ServiceContext) Redis() *redis.Redis {
	if s.zedis == nil {
		s.zedis = s.Config.Redis.NewRedis()
	}
	return s.zedis
}

func (s *ServiceContext) Mysql() *gorm.DB {
	if s.mysql == nil {
		s.mysql = xorm.NewClient(s.Config.Mysql)
	}
	return s.mysql
}

func (s *ServiceContext) ImService() imservice.ImService {
	if s.imService == nil {
		s.imService = imservice.NewImService(zrpc.MustNewClient(s.Config.ImRpc,
			utils.Zrpc.Options()...))
	}
	return s.imService
}

func (s *ServiceContext) RelationService() relationservice.RelationService {
	if s.relationService == nil {
		s.relationService = relationservice.NewRelationService(zrpc.MustNewClient(s.Config.RelationRpc,
			utils.Zrpc.Options()...))
	}
	return s.relationService
}

func (s *ServiceContext) GroupService() groupservice.GroupService {
	if s.groupService == nil {
		s.groupService = groupservice.NewGroupService(zrpc.MustNewClient(s.Config.GroupRpc,
			utils.Zrpc.Options()...))
	}
	return s.groupService
}

func (s *ServiceContext) UserService() userservice.UserService {
	if s.userService == nil {
		s.userService = userservice.NewUserService(zrpc.MustNewClient(s.Config.UserRpc,
			utils.Zrpc.Options()...))
	}
	return s.userService
}

func (s *ServiceContext) NoticeService() noticeservice.NoticeService {
	if s.noticeService == nil {
		s.noticeService = noticeservice.NewNoticeService(zrpc.MustNewClient(s.Config.NoticeRpc,
			utils.Zrpc.Options()...))
	}
	return s.noticeService
}
