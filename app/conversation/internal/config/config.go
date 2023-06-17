package config

import (
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	RedisConf redis.RedisConf
	Group     struct {
		MinGroupId      int `json:",default=100000"`
		MongoCollection xmgo.MongoCollectionConf
	}
	ConversationMember struct {
		MongoCollection xmgo.MongoCollectionConf
	}
	Friend struct {
		MongoCollection xmgo.MongoCollectionConf
	}
	RpcClientConf struct {
		Dispatch zrpc.RpcClientConf
		User     zrpc.RpcClientConf
		Third    zrpc.RpcClientConf
	}
}
