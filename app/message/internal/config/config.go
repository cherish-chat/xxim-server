package config

import (
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	RedisConf redis.RedisConf
	Message   struct {
		MongoCollection xmgo.MongoCollectionConf
	}
	Notice struct {
		MongoCollection xmgo.MongoCollectionConf
	}
}
