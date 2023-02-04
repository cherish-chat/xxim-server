package svc

import (
	"github.com/cherish-chat/xxim-server/app/im/imservice"
	"github.com/cherish-chat/xxim-server/app/notice/noticeservice"
	"github.com/cherish-chat/xxim-server/app/user/internal/config"
	"github.com/cherish-chat/xxim-server/app/user/usermodel"
	"github.com/cherish-chat/xxim-server/common/i18n"
	"github.com/cherish-chat/xxim-server/common/utils/ip2region"
	"github.com/cherish-chat/xxim-server/common/xconf"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config        config.Config
	zedis         *redis.Redis
	mysql         *gorm.DB
	imService     imservice.ImService
	noticeService noticeservice.NoticeService
	ConfigMgr     *xconf.ConfigMgr
	*i18n.I18N
}

func NewServiceContext(c config.Config) *ServiceContext {
	ip2region.Init(c.Ip2RegionUrl)
	s := &ServiceContext{
		Config: c,
	}
	s.ConfigMgr = xconf.NewConfigMgr(s.Mysql(), s.Redis(), "system")
	s.I18N = i18n.NewI18N(s.Mysql())
	usermodel.InitUserSetting(s.Mysql())
	s.Mysql().AutoMigrate(
		&usermodel.User{},
		&usermodel.UserSetting{},
		&usermodel.UserTmp{},
		&usermodel.LoginRecord{},
	)
	return s
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
		s.imService = imservice.NewImService(zrpc.MustNewClient(s.Config.ImRpc))
	}
	return s.imService
}

func (s *ServiceContext) NoticeService() noticeservice.NoticeService {
	if s.noticeService == nil {
		s.noticeService = noticeservice.NewNoticeService(zrpc.MustNewClient(s.Config.NoticeRpc))
	}
	return s.noticeService
}
