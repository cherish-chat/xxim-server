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
	Mysql            xorm.MysqlConfig
	ImRpc            zrpc.RpcClientConf
	RelationRpc      zrpc.RpcClientConf
	GroupRpc         zrpc.RpcClientConf
	UserRpc          zrpc.RpcClientConf
	MobPush          mobpush.Config
	MobAlias         string `json:",default=deviceId,options=deviceId|userId"`
	SyncSendMsgLimit struct {
		Rate  int `json:",default=1"`   // 1 second
		Burst int `json:",default=100"` // 100 qps
	}
}
