package config

import (
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	Mysql                  xorm.MysqlConfig
	ImRpc                  zrpc.RpcClientConf
	NoticeRpc              zrpc.RpcClientConf
	RelationRpc            zrpc.RpcClientConf
	GroupRpc               zrpc.RpcClientConf
	Ip2RegionUrl           string `json:",default=https://github.com/lionsoul2014/ip2region/raw/master/data/ip2region.xdb"`
	EnableMultiDeviceLogin bool   `json:",default=true"`
	Sms                    SmsConfig
}

type SmsConfig struct {
	Enabled    bool
	Type       string
	TencentSms TencentSmsConfig
}

func (c SmsConfig) ToPb() *pb.GetServerAllConfigResp_SmsConfig {
	return &pb.GetServerAllConfigResp_SmsConfig{
		Enabled: c.Enabled,
		Type:    c.Type,
		TencentSms: &pb.GetServerAllConfigResp_TencentSmsConfig{
			AppId:      c.TencentSms.AppId,
			SecretId:   c.TencentSms.SecretId,
			SecretKey:  c.TencentSms.SecretKey,
			Region:     c.TencentSms.Region,
			Sign:       c.TencentSms.Sign,
			TemplateId: c.TencentSms.TemplateId,
		},
	}
}

type TencentSmsConfig struct {
	AppId      string
	SecretId   string
	SecretKey  string
	Region     string
	Sign       string
	TemplateId string
}
