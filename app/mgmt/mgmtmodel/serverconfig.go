package mgmtmodel

import (
	"encoding/json"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"log"
)

var rc *redis.Redis

func InitRedis(redisConfig redis.RedisConf) *redis.Redis {
	rc = redisConfig.NewRedis()
	if !rc.Ping() {
		log.Fatalf("redis init failed, config: %#v", redisConfig)
	}
	// 判断有没有key，没有就创建defaultConfig
	key := rediskey.ServerConfigKey()
	if exist, err := rc.Exists(key); err != nil {
		logx.Errorf("init redis error: %v", err)
		log.Fatalf("redis key %s exists error: %v", key, err)
	} else if !exist {
		config := defaultServerConfig(redisConfig)
		// 存储
		err = rc.Set(key, utils.AnyToString(config))
		if err != nil {
			logx.Errorf("init redis error: %v", err)
			log.Fatalf("redis key %s set error: %v", key, err)
		}
	}
	return rc
}

func GetServerConfigFromRedis() *ServerConfig {
	key := rediskey.ServerConfigKey()
	configString, err := rc.Get(key)
	if err != nil {
		if err == redis.Nil {
			log.Fatalf("redis中没有找到配置信息")
		}
		log.Fatalf("redis中获取配置信息失败: %v", err)
	}
	// 反序列化
	serverConfig := &ServerConfig{}
	err = json.Unmarshal([]byte(configString), serverConfig)
	if err != nil {
		log.Fatalf("反序列化ServerConfig失败: %v", err)
	}
	return serverConfig
}

// ServerConfig 存放在redis中的server配置
type (
	ServerConfig struct {
		Common      ServerCommonConfig
		ConnRpc     ConnRpcConfig
		ImRpc       ImRpcConfig
		MsgRpc      MsgRpcConfig
		UserRpc     UserRpcConfig
		RelationRpc RelationRpcConfig
		GroupRpc    GroupRpcConfig
		NoticeRpc   NoticeRpcConfig
		Mgmt        MgmtConfig
	}
	ServerCommonConfig struct {
		Host         string // default: 0.0.0.0
		RpcTimeOut   int64  // default: 10000
		LogLevel     string // default: info, options: debug,info,error,severe
		Telemetry    TelemetryConfig
		Redis        RedisConfig
		Mysql        xorm.MysqlConfig
		Ip2RegionUrl string // default: https://xxim-public-1312910328.cos.ap-guangzhou.myqcloud.com/ip2region.xdb
		Mode         string // default: pro, options: dev,pro
	}
	TelemetryConfig struct {
		EndPoint string  // default:
		Sampler  float64 // default: 1.0
		Batcher  string  // default: jaeger, options: jaeger|zipkin|grpc
	}
	RedisConfig struct {
		Host string
		Type string // default: node, options=node|cluster"`
		Pass string // default:
		Tls  bool   // default: false
	}
	ConnRpcConfig struct {
		DiscovType    string // default: endpoints, options: k8s|endpoints
		K8sNamespace  string // default: xxim
		Endpoints     []string
		Port          int64 // default: 6700
		WebsocketPort int64 // default: 6701
	}
	ImRpcConfig struct {
		Port int64 // default: 6702
	}
	MsgRpcConfig struct {
		Port    int64 // default: 6703
		MobPush MobPushConfig
		Pulsar  MsgPulsarConfig
	}
	MobPushConfig struct {
		Enabled        bool // default: false
		AppKey         string
		AppSecret      string
		ApnsProduction bool   // default: true
		ApnsCateGory   string // default: ""
		ApnsSound      string // default: "default"
		AndroidSound   string // default: "default"
	}
	MsgPulsarConfig struct {
		Enabled           bool // default: false
		Token             string
		VpcUrl            string
		TopicName         string
		ReceiverQueueSize int64 // default: 1000
		ProducerTimeout   int64 // default: 3000
	}
	UserRpcConfig struct {
		Port int64 // default: 6704
	}
	RelationRpcConfig struct {
		Port int64 // default: 6705
	}
	GroupRpcConfig struct {
		Port int64 // default: 6706
	}
	NoticeRpcConfig struct {
		Port int64 // default: 6707
	}
	MgmtConfig struct {
		RpcPort  int64 // default: 6708
		HttpPort int64 // default: 6799
	}
)

func defaultServerConfig(redisConfig redis.RedisConf) *ServerConfig {
	return &ServerConfig{
		Common: ServerCommonConfig{
			Host:       "0.0.0.0",
			RpcTimeOut: 10000,
			LogLevel:   "info",
			Telemetry:  TelemetryConfig{},
			Redis: RedisConfig{
				Host: redisConfig.Host,
				Type: "node",
				Pass: redisConfig.Pass,
				Tls:  redisConfig.Tls,
			},
			Mysql:        xorm.MysqlConfig{},
			Ip2RegionUrl: "https://xxim-public-1312910328.cos.ap-guangzhou.myqcloud.com/ip2region.xdb",
			Mode:         "pro",
		},
		ConnRpc: ConnRpcConfig{
			DiscovType:   "endpoints",
			K8sNamespace: "xxim",
			Endpoints: []string{
				"127.0.0.1:6700",
			},
			Port:          6700,
			WebsocketPort: 6701,
		},
		ImRpc: ImRpcConfig{
			Port: 6702,
		},
		MsgRpc: MsgRpcConfig{
			Port: 6703,
			MobPush: MobPushConfig{
				Enabled:        false,
				AppKey:         "",
				AppSecret:      "",
				ApnsProduction: true,
				ApnsCateGory:   "",
				ApnsSound:      "default",
				AndroidSound:   "default",
			},
			Pulsar: MsgPulsarConfig{
				Enabled:           false,
				Token:             "",
				VpcUrl:            "",
				TopicName:         "",
				ReceiverQueueSize: 1000,
				ProducerTimeout:   3000,
			},
		},
		UserRpc: UserRpcConfig{
			Port: 6704,
		},
		RelationRpc: RelationRpcConfig{
			Port: 6705,
		},
		GroupRpc: GroupRpcConfig{
			Port: 6706,
		},
		NoticeRpc: NoticeRpcConfig{
			Port: 6707,
		},
		Mgmt: MgmtConfig{
			RpcPort:  6708,
			HttpPort: 6799,
		},
	}
}
