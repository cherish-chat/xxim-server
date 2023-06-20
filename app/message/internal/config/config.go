package config

import (
	"github.com/cherish-chat/xxim-server/common/xmgo"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	RedisConf     redis.RedisConf
	RpcClientConf struct {
		Dispatch     zrpc.RpcClientConf
		User         zrpc.RpcClientConf
		Conversation zrpc.RpcClientConf
		Third        zrpc.RpcClientConf
		Message      zrpc.RpcClientConf
	}
	MongoCollection struct {
		BroadcastNotice           xmgo.MongoCollectionConf
		SubscriptionNotice        xmgo.MongoCollectionConf
		SubscriptionNoticeContent xmgo.MongoCollectionConf
		Message                   xmgo.MongoCollectionConf
	}
	SendMsgLimiter struct {
		Key   string `json:",default=send_msg_limiter"`
		Rate  int    `json:",default=50"`  //每秒钟生成的令牌数
		Burst int    `json:",default=100"` //令牌桶的容量
	}
	//InsertMsgBuffer 插入消息缓冲区
	InsertMsgBuffer struct {
		Size         int `json:",default=1000"` // 缓冲区大小
		LoopInterval int `json:",default=100"`  // 循环间隔 单位(ms)
	}
}
