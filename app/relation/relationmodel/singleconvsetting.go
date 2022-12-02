package relationmodel

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xredis"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

type SingleConvSetting struct {
	ConvId string `gorm:"column:convId;index:conv_user_id,unique;not null;" json:"convId"`
	UserId string `gorm:"column:userId;index:conv_user_id,unique;not null;" json:"userId"`
	// 设为置顶
	IsTop bool `gorm:"column:isTop;default:0;" json:"isTop"`
	// 设为免打扰
	IsDisturb bool `gorm:"column:isDisturb;default:0;" json:"isDisturb"`
	// 消息通知设置 （当免打扰时，此设置无效）
	// 通知显示消息预览
	NotifyPreview bool `gorm:"column:notifyPreview;default:1;" json:"notifyPreview"`
	// 通知声音
	NotifySound bool `gorm:"column:notifySound;default:1;" json:"notifySound"`
	// 通知自定义声音
	NotifyCustomSound string `gorm:"column:notifyCustomSound;default:'';" json:"notifyCustomSound"`
	// 通知震动
	NotifyVibrate bool `gorm:"column:notifyVibrate;default:1;" json:"notifyVibrate"`
	// 屏蔽此人消息
	IsShield bool `gorm:"column:isShield;default:0;" json:"isShield"`
	// 聊天背景
	ChatBg string `gorm:"column:chatBg;default:'';" json:"chatBg"`
}

func (m *SingleConvSetting) TableName() string {
	return "single_conv_settings"
}

func (m *SingleConvSetting) ToProto() *pb.SingleConvSetting {
	return &pb.SingleConvSetting{
		ConvId:            m.ConvId,
		UserId:            m.UserId,
		IsTop:             utils.AnyPtr(m.IsTop),
		IsDisturb:         utils.AnyPtr(m.IsDisturb),
		NotifyPreview:     utils.AnyPtr(m.NotifyPreview),
		NotifySound:       utils.AnyPtr(m.NotifySound),
		NotifyCustomSound: utils.AnyPtr(m.NotifyCustomSound),
		NotifyVibrate:     utils.AnyPtr(m.NotifyVibrate),
		IsShield:          utils.AnyPtr(m.IsShield),
		ChatBg:            utils.AnyPtr(m.ChatBg),
	}
}

func (m *SingleConvSetting) ExpireSeconds() int {
	return xredis.ExpireMinutes(5)
}

func FlushSingleConvSetting(ctx context.Context, rc *redis.Redis, settings ...*SingleConvSetting) error {
	if len(settings) == 0 {
		return nil
	}
	keys := make([]string, 0, len(settings))
	for _, setting := range settings {
		keys = append(keys, rediskey.SingleConvSetting(setting.ConvId, setting.UserId))
	}
	_, err := rc.DelCtx(ctx, keys...)
	return err

}
