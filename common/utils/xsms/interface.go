package xsms

import (
	"errors"
	"github.com/cherish-chat/xxim-server/common/pb"
)

type SmsSender interface {
	SendMsg(mobiles []string, args ...interface{}) error
}

type EmptySmsSender struct {
}

func (e *EmptySmsSender) SendMsg(mobiles []string, args ...interface{}) error {
	return nil
}

func NewSmsSender(config *pb.GetServerAllConfigResp_SmsConfig) (SmsSender, error) {
	if !config.Enabled {
		return &EmptySmsSender{}, nil
	}
	switch config.Type {
	case "tencent":
		smsSender, err := NewTencentSmsSender(config.TencentSms)
		if err != nil {
			return nil, err
		}
		return smsSender, nil
	}
	return nil, errors.New("not support sms type")
}
