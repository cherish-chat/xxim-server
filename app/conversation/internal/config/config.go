package config

import (
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	RedisConf           redis.RedisConf
	ConversationSetting struct {
		MongoCollection xmgo.MongoCollectionConf
	}
	Group struct {
		MinGroupId      int `json:",default=100000"`
		MongoCollection xmgo.MongoCollectionConf
	}
	GroupMember struct {
		MongoCollection xmgo.MongoCollectionConf
	}
	Friend struct {
		MongoCollection xmgo.MongoCollectionConf
	}
}
