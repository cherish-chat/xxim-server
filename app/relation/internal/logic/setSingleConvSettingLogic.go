package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/relation/internal/svc"
	"github.com/cherish-chat/xxim-server/app/relation/relationmodel"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"go.opentelemetry.io/otel/propagation"
	"gorm.io/gorm"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetSingleConvSettingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetSingleConvSettingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetSingleConvSettingLogic {
	return &SetSingleConvSettingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SetSingleConvSettingLogic) SetSingleConvSetting(in *pb.SetSingleConvSettingReq) (*pb.SetSingleConvSettingResp, error) {
	setting := &relationmodel.SingleConvSetting{
		ConvId: in.Setting.ConvId,
		UserId: in.Setting.UserId,
	}
	dest := &relationmodel.SingleConvSetting{}
	err := l.svcCtx.Mysql().Model(setting).Where("convId = ? AND userId = ?", in.Setting.ConvId, in.Setting.UserId).Limit(1).Find(dest).Error
	if err != nil {
		l.Errorf("SetSingleConvSetting: %v", err)
		return &pb.SetSingleConvSettingResp{CommonResp: pb.NewRetryErrorResp()}, nil
	}
	// 删除redis缓存
	err = relationmodel.FlushSingleConvSetting(l.ctx, l.svcCtx.Redis(), setting)
	if err != nil {
		l.Errorf("SetSingleConvSetting: %v", err)
		return &pb.SetSingleConvSettingResp{CommonResp: pb.NewRetryErrorResp()}, nil
	}
	if dest.ConvId == "" {
		// 不存在，插入
		config := l.svcCtx.ConfigMgr.MGetOrDefaultCtx(l.ctx, map[string]string{
			"singleConvSetting_isTop_Default":             "0",
			"singleConvSetting_isDisturb_Default":         "0",
			"singleConvSetting_notifyPreview_Default":     "1",
			"singleConvSetting_notifySound_Default":       "1",
			"singleConvSetting_notifyCustomSound_Default": "",
			"singleConvSetting_notifyVibrate_Default":     "1",
			"singleConvSetting_isShield_Default":          "0",
			"singleConvSetting_chatBg_Default":            "",
		}, in.CommonReq.UserId)
		dest = &relationmodel.SingleConvSetting{
			ConvId:            in.Setting.ConvId,
			UserId:            in.Setting.UserId,
			IsTop:             utils.If(in.Setting.IsTop == nil, utils.String2Bool(config["singleConvSetting_isTop_Default"]), in.Setting.GetIsTop()),
			IsDisturb:         utils.If(in.Setting.IsDisturb == nil, utils.String2Bool(config["singleConvSetting_isDisturb_Default"]), in.Setting.GetIsDisturb()),
			NotifyPreview:     utils.If(in.Setting.NotifyPreview == nil, utils.String2Bool(config["singleConvSetting_notifyPreview_Default"]), in.Setting.GetNotifyPreview()),
			NotifySound:       utils.If(in.Setting.NotifySound == nil, utils.String2Bool(config["singleConvSetting_notifySound_Default"]), in.Setting.GetNotifySound()),
			NotifyCustomSound: utils.If(in.Setting.NotifyCustomSound == nil, config["singleConvSetting_notifyCustomSound_Default"], in.Setting.GetNotifyCustomSound()),
			NotifyVibrate:     utils.If(in.Setting.NotifyVibrate == nil, utils.String2Bool(config["singleConvSetting_notifyVibrate_Default"]), in.Setting.GetNotifyVibrate()),
			IsShield:          utils.If(in.Setting.IsShield == nil, utils.String2Bool(config["singleConvSetting_isShield_Default"]), in.Setting.GetIsShield()),
			ChatBg:            utils.If(in.Setting.ChatBg == nil, config["singleConvSetting_chatBg_Default"], in.Setting.GetChatBg()),
		}
		err = l.svcCtx.Mysql().Model(setting).Create(dest).Error
		if err != nil {
			l.Errorf("SetSingleConvSetting: %v", err)
			return &pb.SetSingleConvSettingResp{CommonResp: pb.NewRetryErrorResp()}, nil
		}
	} else {
		// 存在，更新
		updateMap := make(map[string]interface{})
		if in.Setting.IsTop != nil {
			updateMap["isTop"] = in.Setting.IsTop
		}
		if in.Setting.IsDisturb != nil {
			updateMap["isDisturb"] = in.Setting.IsDisturb
		}
		if in.Setting.NotifyPreview != nil {
			updateMap["notifyPreview"] = in.Setting.NotifyPreview
		}
		if in.Setting.NotifySound != nil {
			updateMap["notifySound"] = in.Setting.NotifySound
		}
		if in.Setting.NotifyCustomSound != nil {
			updateMap["notifyCustomSound"] = in.Setting.NotifyCustomSound
		}
		if in.Setting.NotifyVibrate != nil {
			updateMap["notifyVibrate"] = in.Setting.NotifyVibrate
		}
		if in.Setting.IsShield != nil {
			updateMap["isShield"] = in.Setting.IsShield
		}
		if in.Setting.ChatBg != nil {
			updateMap["chatBg"] = in.Setting.ChatBg
		}
		if len(updateMap) == 0 {
			return &pb.SetSingleConvSettingResp{}, nil
		}
		err := xorm.Transaction(l.svcCtx.Mysql(), func(tx *gorm.DB) error {
			err := xorm.Update(tx, setting, updateMap, xorm.Where("convId = ? AND userId = ?", in.Setting.ConvId, in.Setting.UserId))
			if err != nil {
				l.Errorf("SetSingleConvSetting: %v", err)
			}
			return err
		}, func(tx *gorm.DB) error {
			return nil
		})
		if err != nil {
			l.Errorf("Transaction error: %v", err)
			return &pb.SetSingleConvSettingResp{CommonResp: pb.NewRetryErrorResp()}, nil
		}
	}
	// 缓存预热
	go xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "CacheWarm", func(ctx context.Context) {
		_, err := NewGetSingleConvSettingLogic(ctx, l.svcCtx).GetSingleConvSetting(&pb.GetSingleConvSettingReq{
			CommonReq: in.CommonReq,
			ConvId:    in.Setting.ConvId,
			UserId:    in.Setting.UserId,
		})
		if err != nil {
			l.Errorf("CacheWarm: %v", err)
			return
		}
	}, propagation.MapCarrier{})
	return &pb.SetSingleConvSettingResp{}, nil
}
