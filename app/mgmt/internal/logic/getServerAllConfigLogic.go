package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetServerAllConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetServerAllConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetServerAllConfigLogic {
	return &GetServerAllConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetServerAllConfigLogic) GetServerAllConfig(in *pb.GetServerAllConfigReq) (*pb.GetServerAllConfigResp, error) {
	config := mgmtmodel.GetServerConfigFromRedis()
	return &pb.GetServerAllConfigResp{
		Common: &pb.GetServerAllConfigResp_CommonConfig{
			Host:       config.Common.Host,
			RpcTimeOut: config.Common.RpcTimeOut,
			LogLevel:   config.Common.LogLevel,
			Telemetry: &pb.GetServerAllConfigResp_TelemetryConfig{
				EndPoint: config.Common.Telemetry.EndPoint,
				Sampler:  float32(config.Common.Telemetry.Sampler),
				Batcher:  config.Common.Telemetry.Batcher,
			},
			Redis: &pb.GetServerAllConfigResp_RedisConfig{
				Host: config.Common.Redis.Host,
				Type: config.Common.Redis.Type,
				Pass: config.Common.Redis.Pass,
				Tls:  config.Common.Redis.Tls,
			},
			Mysql: &pb.GetServerAllConfigResp_MysqlConfig{
				Addr:         config.Common.Mysql.Addr,
				MaxIdleConns: int32(config.Common.Mysql.MaxIdleConns),
				MaxOpenConns: int32(config.Common.Mysql.MaxOpenConns),
				LogLevel:     config.Common.Mysql.LogLevel,
			},
			Ip2RegionUrl: config.Common.Ip2RegionUrl,
			Mode:         config.Common.Mode,
		},
		ConnRpc: &pb.GetServerAllConfigResp_ConnRpcConfig{
			DiscovType:    config.ConnRpc.DiscovType,
			K8SNamespace:  config.ConnRpc.K8sNamespace,
			Endpoints:     config.ConnRpc.Endpoints,
			Port:          config.ConnRpc.Port,
			WebsocketPort: config.ConnRpc.WebsocketPort,
		},
		ImRpc: &pb.GetServerAllConfigResp_ImRpcConfig{Port: config.ImRpc.Port},
		MsgRpc: &pb.GetServerAllConfigResp_MsgRpcConfig{
			Port: config.MsgRpc.Port,
			MobPush: &pb.GetServerAllConfigResp_MobPushConfig{
				Enabled:        config.MsgRpc.MobPush.Enabled,
				AppKey:         config.MsgRpc.MobPush.AppKey,
				AppSecret:      config.MsgRpc.MobPush.AppSecret,
				ApnsProduction: config.MsgRpc.MobPush.ApnsProduction,
				ApnsCateGory:   config.MsgRpc.MobPush.ApnsCateGory,
				ApnsSound:      config.MsgRpc.MobPush.ApnsSound,
				AndroidSound:   config.MsgRpc.MobPush.AndroidSound,
			},
			Pulsar: &pb.GetServerAllConfigResp_MsgPulsarConfig{
				Enabled:           config.MsgRpc.Pulsar.Enabled,
				Token:             config.MsgRpc.Pulsar.Token,
				VpcUrl:            config.MsgRpc.Pulsar.VpcUrl,
				TopicName:         config.MsgRpc.Pulsar.TopicName,
				ReceiverQueueSize: config.MsgRpc.Pulsar.ReceiverQueueSize,
				ProducerTimeout:   config.MsgRpc.Pulsar.ProducerTimeout,
			},
			DiscovType:   config.MsgRpc.DiscovType,
			K8SNamespace: config.MsgRpc.K8sNamespace,
			Endpoints:    config.MsgRpc.Endpoints,
		},
		UserRpc:     &pb.GetServerAllConfigResp_UserRpcConfig{Port: config.UserRpc.Port},
		RelationRpc: &pb.GetServerAllConfigResp_RelationRpcConfig{Port: config.RelationRpc.Port},
		GroupRpc: &pb.GetServerAllConfigResp_GroupRpcConfig{
			Port:                config.GroupRpc.Port,
			MaxGroupCount:       config.GroupRpc.MaxGroupCount,
			MaxGroupMemberCount: config.GroupRpc.MaxGroupMemberCount,
		},
		NoticeRpc:  &pb.GetServerAllConfigResp_NoticeRpcConfig{Port: config.NoticeRpc.Port},
		AppMgmtRpc: &pb.GetServerAllConfigResp_AppMgmtRpcConfig{Port: config.AppMgmtRpc.Port},
		Mgmt: &pb.GetServerAllConfigResp_MgmtConfig{
			RpcPort:        config.Mgmt.RpcPort,
			HttpPort:       config.Mgmt.HttpPort,
			SuperAdminId:   config.Mgmt.SuperAdminId,
			SuperAdminPass: config.Mgmt.SuperAdminPass,
		},
	}, nil
}
