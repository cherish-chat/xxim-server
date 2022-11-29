package config

import (
	"github.com/cherish-chat/xxim-server/common/pkg/mobpush"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xtdmq"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	TDMQ struct {
		xtdmq.TDMQConfig
		xtdmq.TDMQConsumerConfig
		Producers struct {
			Msg xtdmq.TDMQProducerConfig
		}
	}
	Mysql       xorm.MysqlConfig
	ImRpc       zrpc.RpcClientConf
	RelationRpc zrpc.RpcClientConf
	GroupRpc    zrpc.RpcClientConf
	MobPush     mobpush.Config
}
