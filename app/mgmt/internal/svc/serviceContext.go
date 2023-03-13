package svc

import (
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtservice"
	"github.com/cherish-chat/xxim-server/app/group/groupservice"
	"github.com/cherish-chat/xxim-server/app/im/imservice"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/config"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"
	msgservice "github.com/cherish-chat/xxim-server/app/msg/msgService"
	"github.com/cherish-chat/xxim-server/app/notice/noticeservice"
	"github.com/cherish-chat/xxim-server/app/relation/relationservice"
	"github.com/cherish-chat/xxim-server/app/user/userservice"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/utils/ip2region"
	"github.com/cherish-chat/xxim-server/common/xconf"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config          config.Config
	redis           *redis.Redis
	mysql           *gorm.DB
	imService       imservice.ImService
	msgService      msgservice.MsgService
	noticeService   noticeservice.NoticeService
	appMgmtService  appmgmtservice.AppMgmtService
	relationService relationservice.RelationService
	userService     userservice.UserService
	groupService    groupservice.GroupService
	ConfigMgr       *xconf.ConfigMgr
}

func NewServiceContext(c config.Config, rc *redis.Redis) *ServiceContext {
	ip2region.Init(c.Ip2RegionUrl)
	s := &ServiceContext{
		Config: c,
		redis:  rc,
	}
	s.ConfigMgr = xconf.NewConfigMgr(s.Mysql(), s.Redis(), "system")
	return s
}

func (s *ServiceContext) Redis() *redis.Redis {
	return s.redis
}

func (s *ServiceContext) Mysql() *gorm.DB {
	if s.mysql == nil {
		s.mysql = xorm.NewClient(s.Config.Mysql)
		s.mysql.AutoMigrate(&appmgmtmodel.Config{})
		s.mysql.AutoMigrate(&mgmtmodel.AutoIncrement{})
		s.mysql.AutoMigrate(&mgmtmodel.User{})
		s.mysql.AutoMigrate(&mgmtmodel.LoginRecord{})
		s.mysql.AutoMigrate(&mgmtmodel.Menu{})
		s.mysql.AutoMigrate(&mgmtmodel.Role{})
		s.mysql.AutoMigrate(&mgmtmodel.ApiPath{})
		s.mysql.AutoMigrate(&mgmtmodel.OperationLog{})
		s.mysql.AutoMigrate(&mgmtmodel.MSIPWhitelist{})
		s.mysql.AutoMigrate(&mgmtmodel.Album{})
		s.mysql.AutoMigrate(&mgmtmodel.AlbumCate{})
		mgmtmodel.InitData(s.mysql)
	}
	return s.mysql
}

func (s *ServiceContext) ImService() imservice.ImService {
	if s.imService == nil {
		s.imService = imservice.NewImService(zrpc.MustNewClient(
			s.Config.ImRpc,
			utils.Zrpc.Options()...))
	}
	return s.imService
}

func (s *ServiceContext) MsgService() msgservice.MsgService {
	if s.msgService == nil {
		s.msgService = msgservice.NewMsgService(zrpc.MustNewClient(s.Config.MsgRpc,
			utils.Zrpc.Options()...))
	}
	return s.msgService
}

func (s *ServiceContext) NoticeService() noticeservice.NoticeService {
	if s.noticeService == nil {
		s.noticeService = noticeservice.NewNoticeService(zrpc.MustNewClient(s.Config.NoticeRpc,
			utils.Zrpc.Options()...))
	}
	return s.noticeService
}

func (s *ServiceContext) AppMgmtService() appmgmtservice.AppMgmtService {
	if s.appMgmtService == nil {
		s.appMgmtService = appmgmtservice.NewAppMgmtService(zrpc.MustNewClient(zrpc.RpcClientConf{
			Etcd:      s.Config.AppMgmtRpc.Etcd,
			Endpoints: s.Config.AppMgmtRpc.Endpoints,
			Target:    s.Config.AppMgmtRpc.Target,
			App:       s.Config.AppMgmtRpc.App,
			Token:     s.Config.AppMgmtRpc.Token,
			NonBlock:  true,
			Timeout:   60 * 1000,
		},
			utils.Zrpc.Options()...))
	}
	return s.appMgmtService
}

func (s *ServiceContext) RelationService() relationservice.RelationService {
	if s.relationService == nil {
		s.relationService = relationservice.NewRelationService(zrpc.MustNewClient(s.Config.RelationRpc,
			utils.Zrpc.Options()...))
	}
	return s.relationService
}

func (s *ServiceContext) UserService() userservice.UserService {
	if s.userService == nil {
		s.userService = userservice.NewUserService(zrpc.MustNewClient(s.Config.UserRpc,
			utils.Zrpc.Options()...))
	}
	return s.userService
}

func (s *ServiceContext) GroupService() groupservice.GroupService {
	if s.groupService == nil {
		s.groupService = groupservice.NewGroupService(zrpc.MustNewClient(s.Config.GroupRpc,
			utils.Zrpc.Options()...))
	}
	return s.groupService
}
