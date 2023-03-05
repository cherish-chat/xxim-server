package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateServerConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateServerConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateServerConfigLogic {
	return &UpdateServerConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateServerConfigLogic) UpdateServerConfig(in *pb.UpdateServerConfigReq) (*pb.UpdateServerConfigResp, error) {
	key := rediskey.ServerConfigKey()
	config := &mgmtmodel.ServerConfig{
		Common: mgmtmodel.ServerCommonConfig{
			Host:       in.Config.Common.Host,
			RpcTimeOut: in.Config.Common.RpcTimeOut,
			LogLevel:   in.Config.Common.LogLevel,
			Telemetry: mgmtmodel.TelemetryConfig{
				EndPoint: in.Config.Common.Telemetry.EndPoint,
				Sampler:  float64(in.Config.Common.Telemetry.Sampler),
				Batcher:  in.Config.Common.Telemetry.Batcher,
			},
			Redis: mgmtmodel.RedisConfig{
				Host: in.Config.Common.Redis.Host,
				Type: in.Config.Common.Redis.Type,
				Pass: in.Config.Common.Redis.Pass,
				Tls:  in.Config.Common.Redis.Tls,
			},
			Mysql: xorm.MysqlConfig{
				Addr:         in.Config.Common.Mysql.Addr,
				MaxIdleConns: int(in.Config.Common.Mysql.MaxIdleConns),
				MaxOpenConns: int(in.Config.Common.Mysql.MaxOpenConns),
				LogLevel:     in.Config.Common.Mysql.LogLevel,
			},
			Ip2RegionUrl: in.Config.Common.Ip2RegionUrl,
			Mode:         in.Config.Common.Mode,
		},
		ConnRpc: mgmtmodel.ConnRpcConfig{
			DiscovType:    in.Config.ConnRpc.DiscovType,
			K8sNamespace:  in.Config.ConnRpc.K8SNamespace,
			Endpoints:     in.Config.ConnRpc.Endpoints,
			Port:          in.Config.ConnRpc.Port,
			WebsocketPort: in.Config.ConnRpc.WebsocketPort,
			RsaPublicKey:  in.Config.ConnRpc.RsaPublicKey,
			RsaPrivateKey: in.Config.ConnRpc.RsaPrivateKey,
		},
		ImRpc: mgmtmodel.ImRpcConfig{Port: in.Config.ImRpc.Port},
		MsgRpc: mgmtmodel.MsgRpcConfig{
			DiscovType:   in.Config.MsgRpc.DiscovType,
			K8sNamespace: in.Config.MsgRpc.K8SNamespace,
			Endpoints:    in.Config.MsgRpc.Endpoints,
			Port:         in.Config.MsgRpc.Port,
			MobPush: mgmtmodel.MobPushConfig{
				Enabled:        in.Config.MsgRpc.MobPush.Enabled,
				AppKey:         in.Config.MsgRpc.MobPush.AppKey,
				AppSecret:      in.Config.MsgRpc.MobPush.AppSecret,
				ApnsProduction: in.Config.MsgRpc.MobPush.ApnsProduction,
				ApnsCateGory:   in.Config.MsgRpc.MobPush.ApnsCateGory,
				ApnsSound:      in.Config.MsgRpc.MobPush.ApnsSound,
				AndroidSound:   in.Config.MsgRpc.MobPush.AndroidSound,
			},
			Pulsar: mgmtmodel.MsgPulsarConfig{
				Enabled:           in.Config.MsgRpc.Pulsar.Enabled,
				Token:             in.Config.MsgRpc.Pulsar.Token,
				VpcUrl:            in.Config.MsgRpc.Pulsar.VpcUrl,
				TopicName:         in.Config.MsgRpc.Pulsar.TopicName,
				ReceiverQueueSize: in.Config.MsgRpc.Pulsar.ReceiverQueueSize,
				ProducerTimeout:   in.Config.MsgRpc.Pulsar.ProducerTimeout,
			},
		},
		UserRpc: mgmtmodel.UserRpcConfig{Port: in.Config.UserRpc.Port, Sms: mgmtmodel.SmsConfig{
			Enabled: in.Config.UserRpc.Sms.Enabled,
			Type:    in.Config.UserRpc.Sms.Type,
			TencentSms: mgmtmodel.TencentSmsConfig{
				AppId:      in.Config.UserRpc.Sms.TencentSms.AppId,
				SecretId:   in.Config.UserRpc.Sms.TencentSms.SecretId,
				SecretKey:  in.Config.UserRpc.Sms.TencentSms.SecretKey,
				Region:     in.Config.UserRpc.Sms.TencentSms.Region,
				Sign:       in.Config.UserRpc.Sms.TencentSms.Sign,
				TemplateId: in.Config.UserRpc.Sms.TencentSms.TemplateId,
			},
		}},
		RelationRpc: mgmtmodel.RelationRpcConfig{Port: in.Config.RelationRpc.Port},
		GroupRpc: mgmtmodel.GroupRpcConfig{
			Port:                in.Config.GroupRpc.Port,
			MaxGroupCount:       in.Config.GroupRpc.MaxGroupCount,
			MaxGroupMemberCount: in.Config.GroupRpc.MaxGroupMemberCount,
		},
		NoticeRpc:  mgmtmodel.NoticeRpcConfig{Port: in.Config.NoticeRpc.Port},
		AppMgmtRpc: mgmtmodel.AppMgmtRpcConfig{Port: in.Config.AppMgmtRpc.Port},
		Mgmt: mgmtmodel.MgmtConfig{
			RpcPort:        in.Config.Mgmt.RpcPort,
			HttpPort:       in.Config.Mgmt.HttpPort,
			SuperAdminId:   in.Config.Mgmt.SuperAdminId,
			SuperAdminPass: in.Config.Mgmt.SuperAdminPass,
		},
	}
	err := l.svcCtx.Redis().SetCtx(l.ctx, key, utils.AnyToString(config))
	if err != nil {
		l.Errorf("set server config to redis error: %v", err)
		return &pb.UpdateServerConfigResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.UpdateServerConfigResp{}, nil
}
