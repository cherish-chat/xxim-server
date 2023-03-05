package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/im/immodel"
	"github.com/cherish-chat/xxim-server/app/notice/noticemodel"
	"github.com/cherish-chat/xxim-server/common/utils"
	"github.com/cherish-chat/xxim-server/common/xorm"
	"github.com/cherish-chat/xxim-server/common/xtrace"
	"go.opentelemetry.io/otel/propagation"
	"gorm.io/gorm"
	"time"

	"github.com/cherish-chat/xxim-server/app/im/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateConvSettingLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateConvSettingLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateConvSettingLogic {
	return &UpdateConvSettingLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateConvSettingLogic) UpdateConvSetting(in *pb.UpdateConvSettingReq) (*pb.UpdateConvSettingResp, error) {
	setting := &immodel.ConvSetting{
		ConvId: in.ConvSetting.ConvId,
		UserId: in.CommonReq.UserId,
	}
	dest := &immodel.ConvSetting{}
	err := l.svcCtx.Mysql().Model(setting).Where("convId = ? AND userId = ?", in.ConvSetting.ConvId, in.CommonReq.UserId).Limit(1).Find(dest).Error
	if err != nil {
		l.Errorf("SetConvSetting: %v", err)
		return &pb.UpdateConvSettingResp{CommonResp: pb.NewRetryErrorResp()}, nil
	}
	// 删除redis缓存
	err = immodel.FlushConvSetting(l.ctx, l.svcCtx.Redis(), setting)
	if err != nil {
		l.Errorf("SetConvSetting: %v", err)
		return &pb.UpdateConvSettingResp{CommonResp: pb.NewRetryErrorResp()}, nil
	}
	err = xorm.Transaction(l.svcCtx.Mysql(), func(tx *gorm.DB) error {
		if dest.ConvId == "" {
			// 不存在，插入
			dest = immodel.DefaultConvSetting(in.CommonReq.UserId, in.ConvSetting.ConvId)
			if in.ConvSetting.IsTop != nil {
				dest.IsTop = *in.ConvSetting.IsTop
			}
			if in.ConvSetting.IsDisturb != nil {
				dest.IsDisturb = *in.ConvSetting.IsDisturb
			}
			if in.ConvSetting.NotifyPreview != nil {
				dest.NotifyPreview = *in.ConvSetting.NotifyPreview
			}
			if in.ConvSetting.NotifySound != nil {
				dest.NotifySound = *in.ConvSetting.NotifySound
			}
			if in.ConvSetting.NotifyCustomSound != nil {
				dest.NotifyCustomSound = *in.ConvSetting.NotifyCustomSound
			}
			if in.ConvSetting.NotifyVibrate != nil {
				dest.NotifyVibrate = *in.ConvSetting.NotifyVibrate
			}
			if in.ConvSetting.IsShield != nil {
				dest.IsShield = *in.ConvSetting.IsShield
			}
			if in.ConvSetting.ChatBg != nil {
				dest.ChatBg = *in.ConvSetting.ChatBg
			}
			err = tx.Model(setting).Create(dest).Error
			if err != nil {
				l.Errorf("SetConvSetting: %v", err)
			}
			return err
		} else {
			// 存在，更新
			updateMap := make(map[string]interface{})
			if in.ConvSetting.IsTop != nil {
				updateMap["isTop"] = in.ConvSetting.IsTop
			}
			if in.ConvSetting.IsDisturb != nil {
				updateMap["isDisturb"] = in.ConvSetting.IsDisturb
			}
			if in.ConvSetting.NotifyPreview != nil {
				updateMap["notifyPreview"] = in.ConvSetting.NotifyPreview
			}
			if in.ConvSetting.NotifySound != nil {
				updateMap["notifySound"] = in.ConvSetting.NotifySound
			}
			if in.ConvSetting.NotifyCustomSound != nil {
				updateMap["notifyCustomSound"] = in.ConvSetting.NotifyCustomSound
			}
			if in.ConvSetting.NotifyVibrate != nil {
				updateMap["notifyVibrate"] = in.ConvSetting.NotifyVibrate
			}
			if in.ConvSetting.IsShield != nil {
				updateMap["isShield"] = in.ConvSetting.IsShield
			}
			if in.ConvSetting.ChatBg != nil {
				updateMap["chatBg"] = in.ConvSetting.ChatBg
			}
			if len(updateMap) == 0 {
				return nil
			}
			err := xorm.Update(tx, setting, updateMap, xorm.Where("convId = ? AND userId = ?", in.ConvSetting.ConvId, in.CommonReq.UserId))
			if err != nil {
				l.Errorf("SetConvSetting: %v", err)
			}
			return err
		}
	}, func(tx *gorm.DB) error {
		notice := &noticemodel.Notice{
			ConvId: pb.HiddenConvId(in.ConvSetting.ConvId),
			Options: noticemodel.NoticeOption{
				StorageForClient: false,
				UpdateConvNotice: false,
			},
			ContentType: pb.NoticeContentType_SyncConvSetting,
			Content: utils.AnyToBytes(pb.NoticeContent_SyncConvSetting{
				ConvIds: []string{in.ConvSetting.ConvId},
				UserId:  in.CommonReq.UserId,
			}),
			UniqueId: in.CommonReq.UserId,
			UserId:   in.CommonReq.UserId,
			Title:    "",
			Ext:      nil,
		}
		err = notice.Insert(l.ctx, tx, l.svcCtx.Redis())
		if err != nil {
			l.Errorf("insert notice failed, err: %v", err)
		}
		return err
	}, func(tx *gorm.DB) error {
		return immodel.FlushConvSetting(l.ctx, l.svcCtx.Redis(), setting)
	})
	if err != nil {
		return &pb.UpdateConvSettingResp{CommonResp: pb.NewRetryErrorResp()}, nil
	}
	// 缓存预热
	go xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "CacheWarm", func(ctx context.Context) {
		_, err := NewGetConvSettingLogic(ctx, l.svcCtx).GetConvSetting(&pb.GetConvSettingReq{
			CommonReq: in.CommonReq,
			ConvIds:   []string{in.ConvSetting.ConvId},
		})
		if err != nil {
			l.Errorf("CacheWarm: %v", err)
			return
		}
	}, propagation.MapCarrier{})
	// 通知客户端
	utils.RetryProxy(l.ctx, 12, time.Second, func() error {
		_, err = l.svcCtx.NoticeService().GetUserNoticeData(l.ctx, &pb.GetUserNoticeDataReq{
			CommonReq: in.CommonReq,
			UserId:    in.CommonReq.UserId,
			ConvId:    pb.HiddenConvId(in.ConvSetting.ConvId),
			DeviceId:  nil,
		})
		if err != nil {
			l.Errorf("GetUserNoticeData: %v", err)
		}
		return err
	})
	return &pb.UpdateConvSettingResp{}, nil
}
