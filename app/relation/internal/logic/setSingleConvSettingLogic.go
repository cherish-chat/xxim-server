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
		config := l.svcCtx.SystemConfigMgr.MGetOrDefaultCtx(l.ctx, map[string]string{
			"singleConvSetting_isTop_Default":             "0",
			"singleConvSetting_isDisturb_Default":         "0",
			"singleConvSetting_notifyPreview_Default":     "1",
			"singleConvSetting_notifySound_Default":       "1",
			"singleConvSetting_notifyCustomSound_Default": "",
			"singleConvSetting_notifyVibrate_Default":     "1",
			"singleConvSetting_isShield_Default":          "0",
			"singleConvSetting_chatBg_Default":            "",
		})
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
		err := xorm.Update(l.svcCtx.Mysql(), setting, updateMap, xorm.Where("convId = ? AND userId = ?", in.Setting.ConvId, in.Setting.UserId))
		if err != nil {
			l.Errorf("SetSingleConvSetting: %v", err)
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
		// TODO 发送消息通知
		//xtrace.StartFuncSpan(ctx, "SendMsg", func(ctx context.Context) {
		//	bytes, _ := proto.Marshal(resp)
		//	for i := 0; i < 10; i++ {
		//		_, err := msgservice.SendMsgSync(l.svcCtx.MsgService(), ctx, []*pb.MsgData{
		//			msgmodel.CreateConvProfileChangeMsg(in.Setting.UserId, in.Setting.ConvId, bytes).ToMsgData(),
		//		})
		//		if err != nil {
		//			l.Errorf("SendMsg: %v", err)
		//			time.Sleep(time.Second)
		//			continue
		//		}
		//		break
		//	}
		//})
	}, propagation.MapCarrier{})
	return &pb.SetSingleConvSettingResp{}, nil
}
