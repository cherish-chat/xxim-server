package logic

import (
	"context"
	"encoding/json"
	"github.com/cherish-chat/xxim-server/app/relation/relationmodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xredis"
	"github.com/cherish-chat/xxim-server/common/xredis/rediskey"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"go.opentelemetry.io/otel/propagation"

	"github.com/cherish-chat/xxim-server/app/relation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSingleConvSettingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSingleConvSettingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSingleConvSettingLogic {
	return &GetSingleConvSettingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSingleConvSettingLogic) GetSingleConvSetting(in *pb.GetSingleConvSettingReq) (*pb.GetSingleConvSettingResp, error) {
	// 从缓存获取
	setting := &relationmodel.SingleConvSetting{
		ConvId: in.ConvId,
		UserId: in.UserId,
	}
	val, err := l.svcCtx.Redis().GetCtx(l.ctx, rediskey.SingleConvSetting(in.ConvId, in.UserId))
	if err != nil || val == "" {
		return l.getSingleConvSettingFromDB(in)
	}
	if val == xredis.NotFound {
		// 不存在，插入
		return l.notFound(in)
	}
	err = json.Unmarshal([]byte(val), setting)
	if err != nil {
		l.Errorf("SetSingleConvSetting: %v", err)
		return l.getSingleConvSettingFromDB(in)
	}
	return &pb.GetSingleConvSettingResp{
		CommonResp: pb.NewSuccessResp(),
		Setting:    setting.ToProto(),
	}, nil
}

func (l *GetSingleConvSettingLogic) getSingleConvSettingFromDB(in *pb.GetSingleConvSettingReq) (*pb.GetSingleConvSettingResp, error) {
	// 从数据库获取
	setting := &relationmodel.SingleConvSetting{
		ConvId: in.ConvId,
		UserId: in.UserId,
	}
	err := l.svcCtx.Mysql().Model(setting).Where("convId = ? and userId = ?", in.ConvId, in.UserId).First(setting).Error
	if err != nil {
		if xorm.RecordNotFound(err) {
			// 不存在，插入
			return l.notFound(in)
		}
		l.Errorf("SetSingleConvSetting: %v", err)
		return &pb.GetSingleConvSettingResp{CommonResp: pb.NewRetryErrorResp()}, nil
	}
	// 设置缓存
	go xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "SetCache", func(ctx context.Context) {
		err = l.svcCtx.Redis().SetexCtx(ctx, rediskey.SingleConvSetting(in.ConvId, in.UserId), utils.AnyToString(setting), setting.ExpireSeconds())
		if err != nil {
			l.Errorf("SetCache: %v", err)
		}
	}, propagation.MapCarrier{})
	return &pb.GetSingleConvSettingResp{
		Setting: setting.ToProto(),
	}, nil
}

func (l *GetSingleConvSettingLogic) notFound(in *pb.GetSingleConvSettingReq) (*pb.GetSingleConvSettingResp, error) {
	config := l.svcCtx.ConfigMgr.MGetOrDefaultCtx(l.ctx, map[string]string{
		"singleConvSetting_isTop_Default":             "0",
		"singleConvSetting_isDisturb_Default":         "0",
		"singleConvSetting_notifyPreview_Default":     "1",
		"singleConvSetting_notifySound_Default":       "1",
		"singleConvSetting_notifyCustomSound_Default": "",
		"singleConvSetting_notifyVibrate_Default":     "1",
		"singleConvSetting_isShield_Default":          "0",
		"singleConvSetting_chatBg_Default":            "",
	})
	dest := &relationmodel.SingleConvSetting{
		ConvId:            in.ConvId,
		UserId:            in.UserId,
		IsTop:             utils.String2Bool(config["singleConvSetting_isTop_Default"]),
		IsDisturb:         utils.String2Bool(config["singleConvSetting_isDisturb_Default"]),
		NotifyPreview:     utils.String2Bool(config["singleConvSetting_notifyPreview_Default"]),
		NotifySound:       utils.String2Bool(config["singleConvSetting_notifySound_Default"]),
		NotifyCustomSound: config["singleConvSetting_notifyCustomSound_Default"],
		NotifyVibrate:     utils.String2Bool(config["singleConvSetting_notifyVibrate_Default"]),
		IsShield:          utils.String2Bool(config["singleConvSetting_isShield_Default"]),
		ChatBg:            config["singleConvSetting_chatBg_Default"],
	}
	err := relationmodel.FlushSingleConvSetting(l.ctx, l.svcCtx.Redis(), dest)
	if err != nil {
		l.Errorf("SetSingleConvSetting: %v", err)
		return &pb.GetSingleConvSettingResp{CommonResp: pb.NewRetryErrorResp()}, nil
	}
	err = l.svcCtx.Mysql().Model(dest).Create(dest).Error
	if err != nil {
		l.Errorf("SetSingleConvSetting: %v", err)
		return &pb.GetSingleConvSettingResp{CommonResp: pb.NewRetryErrorResp()}, nil
	}
	return &pb.GetSingleConvSettingResp{
		Setting: dest.ToProto(),
	}, nil
}
