package immodel

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xredis"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/gorm"
)

type ConvSetting struct {
	ConvId string `gorm:"column:convId;index:conv_user_id,unique;not null;index;" json:"convId"`
	UserId string `gorm:"column:userId;index:conv_user_id,unique;not null;" json:"userId"`
	// 设为置顶
	IsTop bool `gorm:"column:isTop;default:0;" json:"isTop"`
	// 设为免打扰
	IsDisturb bool `gorm:"column:isDisturb;default:0;index;" json:"isDisturb"`
	// 消息通知设置 （当免打扰时，此设置无效）
	// 通知显示消息预览
	NotifyPreview bool `gorm:"column:notifyPreview;default:1;" json:"notifyPreview"`
	// 通知声音
	NotifySound bool `gorm:"column:notifySound;default:1;" json:"notifySound"`
	// 通知自定义声音
	NotifyCustomSound string `gorm:"column:notifyCustomSound;default:'';" json:"notifyCustomSound"`
	// 通知震动
	NotifyVibrate bool `gorm:"column:notifyVibrate;default:1;" json:"notifyVibrate"`
	// 屏蔽消息
	IsShield bool `gorm:"column:isShield;default:0;" json:"isShield"`
	// 聊天背景
	ChatBg string `gorm:"column:chatBg;default:'';" json:"chatBg"`
}

func SearchGroupMemberList(tx *gorm.DB, convId string, limit int, filter map[string]any) ([]*ConvSetting, error) {
	// TODO: 使用redis优化
	var convSettings []*ConvSetting
	tx = tx.Model(&ConvSetting{}).Where("convId = ?", convId)
	for k, v := range filter {
		tx = tx.Where(k+" = ?", v)
	}
	err := tx.Limit(limit).Find(&convSettings).Error
	return convSettings, err
}

func DefaultConvSetting(userId string, convId string) *ConvSetting {
	return &ConvSetting{
		ConvId:            convId,
		UserId:            userId,
		IsTop:             false,
		IsDisturb:         true,
		NotifyPreview:     true,
		NotifySound:       true,
		NotifyCustomSound: "",
		NotifyVibrate:     true,
		IsShield:          false,
		ChatBg:            "",
	}
}

func (m *ConvSetting) TableName() string {
	return "conv_settings"
}

func (m *ConvSetting) ToProto() *pb.ConvSetting {
	return &pb.ConvSetting{
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

func (m *ConvSetting) ExpireSeconds() int {
	return xredis.ExpireMinutes(5)
}

func FlushConvSetting(ctx context.Context, rc *redis.Redis, settings ...*ConvSetting) error {
	if len(settings) == 0 {
		return nil
	}
	keys := make([]string, 0, len(settings))
	for _, setting := range settings {
		keys = append(keys, rediskey.ConvSetting(setting.ConvId, setting.UserId))
	}
	_, err := rc.DelCtx(ctx, keys...)
	return err

}
