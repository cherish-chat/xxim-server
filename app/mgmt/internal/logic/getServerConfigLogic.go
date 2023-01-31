package logic

import (
	"context"
	"fmt"
	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetServerConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetServerConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetServerConfigLogic {
	return &GetServerConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetServerConfigLogic) GetServerConfig(in *pb.GetServerConfigReq) (*pb.GetServerConfigResp, error) {
	c := GetConfig(mgmtmodel.GetServerConfigFromRedis(), in.ServerName)
	// 序列化为yaml
	s := utils.AnyToStringBeautiful(c)
	// 打印配置信息
	logx.Debugf("ServerConfig: %s", s)
	return &pb.GetServerConfigResp{Config: []byte(s)}, nil
}

func MgmtConfig() []byte {
	c := GetConfig(mgmtmodel.GetServerConfigFromRedis(), "mgmt")
	// 序列化为yaml
	bytes := utils.AnyToStringBeautiful(c)
	return []byte(bytes)
}

func GetConfig(serverConfig *mgmtmodel.ServerConfig, name string) map[string]any {
	c := map[string]any{
		"Name":         name,
		"Timeout":      serverConfig.Common.RpcTimeOut,
		"Mode":         serverConfig.Common.Mode,
		"Ip2RegionUrl": serverConfig.Common.Ip2RegionUrl,

		"Log": map[string]any{
			"ServiceName":         name,
			"Mode":                "console",
			"Encoding":            "json",
			"TimeFormat":          "2006-01-02_15:04:05.000",
			"Path":                "logs",
			"Level":               serverConfig.Common.LogLevel,
			"Compress":            false,
			"KeepDays":            0,
			"StackCooldownMillis": 100,
			"MaxBackups":          0,
			"MaxSize":             0,
			"Rotation":            "daily",
		},
		"Prometheus": map[string]any{
			"Host": "0.0.0.0",
			"Port": 9101,
			"Path": "/metrics",
		},
		"Telemetry": map[string]any{
			"Name":     name,
			"Endpoint": serverConfig.Common.Telemetry.EndPoint,
			"Sampler":  serverConfig.Common.Telemetry.Sampler,
			"Batcher":  serverConfig.Common.Telemetry.Batcher,
		},
		"Redis": map[string]any{
			"Host": serverConfig.Common.Redis.Host,
			"Type": serverConfig.Common.Redis.Type,
			"Pass": serverConfig.Common.Redis.Pass,
			"Tls":  serverConfig.Common.Redis.Tls,
		},
		"Mysql": map[string]any{
			"Addr":         serverConfig.Common.Mysql.Addr,
			"MaxIdleConns": serverConfig.Common.Mysql.MaxIdleConns,
			"MaxOpenConns": serverConfig.Common.Mysql.MaxOpenConns,
			"LogLevel":     serverConfig.Common.Mysql.LogLevel,
		},
		"ConnRpc": map[string]any{
			"DiscovType": serverConfig.ConnRpc.DiscovType,
			"Endpoints":  serverConfig.ConnRpc.Endpoints,
			"K8s": map[string]any{
				"Namespace": serverConfig.ConnRpc.K8sNamespace,
				"Port":      serverConfig.ConnRpc.Port,
			},
		},
		"ImRpc": map[string]any{
			"Endpoints": []string{fmt.Sprintf("127.0.0.1:%d", serverConfig.ImRpc.Port)},
			"NonBlock":  true,
		},
		"MsgRpc": map[string]any{
			"Endpoints": []string{fmt.Sprintf("127.0.0.1:%d", serverConfig.MsgRpc.Port)},
			"NonBlock":  true,
		},
		"UserRpc": map[string]any{
			"Endpoints": []string{fmt.Sprintf("127.0.0.1:%d", serverConfig.UserRpc.Port)},
			"NonBlock":  true,
		},
		"RelationRpc": map[string]any{
			"Endpoints": []string{fmt.Sprintf("127.0.0.1:%d", serverConfig.RelationRpc.Port)},
			"NonBlock":  true,
		},
		"GroupRpc": map[string]any{
			"Endpoints": []string{fmt.Sprintf("127.0.0.1:%d", serverConfig.GroupRpc.Port)},
			"NonBlock":  true,
		},
		"NoticeRpc": map[string]any{
			"Endpoints": []string{fmt.Sprintf("127.0.0.1:%d", serverConfig.NoticeRpc.Port)},
			"NonBlock":  true,
		},
		"MgmtRpc": map[string]any{
			"Endpoints": []string{fmt.Sprintf("127.0.0.1:%d", serverConfig.Mgmt.RpcPort)},
			"NonBlock":  true,
		},
	}
	switch name {
	case "mgmt":
		c["Gin"] = map[string]any{
			"Cors": map[string]any{
				"AllowOrigins":     []string{"*"},
				"AllowHeaders":     []string{"*"},
				"AllowMethods":     []string{"*"},
				"ExposeHeaders":    []string{"*"},
				"AllowCredentials": true,
			},
			"Addr": fmt.Sprintf("%s:%d", serverConfig.Common.Host, serverConfig.Mgmt.HttpPort),
		}
		c["ListenOn"] = fmt.Sprintf("%s:%d", serverConfig.Common.Host, serverConfig.Mgmt.RpcPort)
	case "conn":
		c["Websocket"] = map[string]any{
			"Port": serverConfig.ConnRpc.WebsocketPort,
			"Host": serverConfig.Common.Host,
		}
		c["ListenOn"] = fmt.Sprintf("%s:%d", serverConfig.Common.Host, serverConfig.ConnRpc.Port)
	case "im":
		c["ListenOn"] = fmt.Sprintf("%s:%d", serverConfig.Common.Host, serverConfig.ImRpc.Port)
	case "group":
		c["ListenOn"] = fmt.Sprintf("%s:%d", serverConfig.Common.Host, serverConfig.GroupRpc.Port)
	case "notice":
		c["ListenOn"] = fmt.Sprintf("%s:%d", serverConfig.Common.Host, serverConfig.NoticeRpc.Port)
	case "relation":
		c["ListenOn"] = fmt.Sprintf("%s:%d", serverConfig.Common.Host, serverConfig.RelationRpc.Port)
	case "user":
		c["ListenOn"] = fmt.Sprintf("%s:%d", serverConfig.Common.Host, serverConfig.UserRpc.Port)
	case "msg":
		c["ListenOn"] = fmt.Sprintf("%s:%d", serverConfig.Common.Host, serverConfig.MsgRpc.Port)
		c["MobPush"] = map[string]any{
			"Enabled":        serverConfig.MsgRpc.MobPush.Enabled,
			"AppKey":         serverConfig.MsgRpc.MobPush.AppKey,
			"AppSecret":      serverConfig.MsgRpc.MobPush.AppSecret,
			"ApnsProduction": serverConfig.MsgRpc.MobPush.ApnsProduction,
			"ApnsCateGory":   serverConfig.MsgRpc.MobPush.ApnsCateGory,
			"ApnsSound":      serverConfig.MsgRpc.MobPush.ApnsSound,
			"AndroidSound":   serverConfig.MsgRpc.MobPush.AndroidSound,
		}
		c["TDMQ"] = map[string]any{
			"Enabled":           serverConfig.MsgRpc.Pulsar.Enabled,
			"Token":             serverConfig.MsgRpc.Pulsar.Token,
			"VpcUrl":            serverConfig.MsgRpc.Pulsar.VpcUrl,
			"TopicName":         serverConfig.MsgRpc.Pulsar.TopicName,
			"ReceiverQueueSize": serverConfig.MsgRpc.Pulsar.ReceiverQueueSize,
			"Producer": map[string]any{
				"TopicName":   serverConfig.MsgRpc.Pulsar.TopicName,
				"SendTimeout": serverConfig.MsgRpc.Pulsar.ProducerTimeout,
			},
		}
	}
	return c
}
