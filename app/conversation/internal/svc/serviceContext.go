package svc

import (
	"github.com/cherish-chat/xxim-server/app/conversation/conversationmodel"
	"github.com/cherish-chat/xxim-server/app/conversation/friendmodel"
	"github.com/cherish-chat/xxim-server/app/conversation/groupmodel"
	"github.com/cherish-chat/xxim-server/app/conversation/internal/config"
	"github.com/cherish-chat/xxim-server/common/xcache"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	Redis  *redis.Redis
	Mysql  *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	s := &ServiceContext{
		Config: c,
		Redis:  xcache.MustNewRedis(c.RedisConf),
		Mysql:  xorm.MustNewMysql(c.MysqlConf),
	}
	groupmodel.InitGroupModel(s.Mysql, s.Redis, s.Config.Group.MinGroupId)
	groupmodel.InitGroupMemberModel(s.Mysql, s.Redis)

	friendmodel.InitFriendModel(s.Mysql, s.Redis)

	conversationmodel.InitConversationSettingModel(s.Mysql, s.Redis)
	return s
}
