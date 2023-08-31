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
	convId := ""
	if in.GetConvSetting().GetConvId() != "" {
		convId = in.GetConvSetting().GetConvId()
	} else if in.GetConvSetting2().GetConvId() != "" {
		convId = in.GetConvSetting2().GetConvId()
	}
	setting := &immodel.ConvSetting{
		ConvId: convId,
		UserId: in.CommonReq.UserId,
	}
	dest := &immodel.ConvSetting{}
	err := l.svcCtx.Mysql().Model(setting).Where("convId = ? AND userId = ?", convId, in.CommonReq.UserId).Limit(1).Find(dest).Error
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
			if in.GetConvSetting().GetConvId() != "" {
				dest = immodel.DefaultConvSetting(in.CommonReq.UserId, convId)
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
			} else if in.GetConvSetting2().GetConvId() != "" {
				dest = immodel.DefaultConvSetting(in.CommonReq.UserId, convId)
				if in.ConvSetting2.IsTop == 1 {
					dest.IsTop = true
				} else if in.ConvSetting2.IsTop == 2 {
					dest.IsTop = false
				}
				if in.ConvSetting2.IsDisturb == 1 {
					dest.IsDisturb = true
				} else if in.ConvSetting2.IsDisturb == 2 {
					dest.IsDisturb = false
				}
				if in.ConvSetting2.NotifyPreview == 1 {
					dest.NotifyPreview = true
				} else if in.ConvSetting2.NotifyPreview == 2 {
					dest.NotifyPreview = false
				}
				if in.ConvSetting2.NotifySound == 1 {
					dest.NotifySound = true
				} else if in.ConvSetting2.NotifySound == 2 {
					dest.NotifySound = false
				}
				dest.NotifyCustomSound = in.ConvSetting2.NotifyCustomSound
				if in.ConvSetting2.NotifyVibrate == 1 {
					dest.NotifyVibrate = true
				} else if in.ConvSetting2.NotifyVibrate == 2 {
					dest.NotifyVibrate = false
				}
				if in.ConvSetting2.IsShield == 1 {
					dest.IsShield = true
				} else if in.ConvSetting2.IsShield == 2 {
					dest.IsShield = false
				}
				dest.ChatBg = in.ConvSetting2.ChatBg
				err = tx.Model(setting).Create(dest).Error
				if err != nil {
					l.Errorf("SetConvSetting: %v", err)
				}
				return err
			} else {
				return nil
			}
		} else {
			// 存在，更新
			updateMap := make(map[string]interface{})
			if in.GetConvSetting().GetConvId() != "" {
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
			} else if in.GetConvSetting2().GetConvId() != "" {
				if in.ConvSetting2.IsTop == 1 {
					updateMap["isTop"] = true
				} else if in.ConvSetting2.IsTop == 2 {
					updateMap["isTop"] = false
				}
				if in.ConvSetting2.IsDisturb == 1 {
					updateMap["isDisturb"] = true
				} else if in.ConvSetting2.IsDisturb == 2 {
					updateMap["isDisturb"] = false
				}
				if in.ConvSetting2.NotifyPreview == 1 {
					updateMap["notifyPreview"] = true
				} else if in.ConvSetting2.NotifyPreview == 2 {
					updateMap["notifyPreview"] = false
				}
				if in.ConvSetting2.NotifySound == 1 {
					updateMap["notifySound"] = true
				} else if in.ConvSetting2.NotifySound == 2 {
					updateMap["notifySound"] = false
				}
				updateMap["notifyCustomSound"] = in.ConvSetting2.NotifyCustomSound
				if in.ConvSetting2.NotifyVibrate == 1 {
					updateMap["notifyVibrate"] = true
				} else if in.ConvSetting2.NotifyVibrate == 2 {
					updateMap["notifyVibrate"] = false
				}
				if in.ConvSetting2.IsShield == 1 {
					updateMap["isShield"] = true
				} else if in.ConvSetting2.IsShield == 2 {
					updateMap["isShield"] = false
				}
				updateMap["chatBg"] = in.ConvSetting2.ChatBg
			}
			if len(updateMap) == 0 {
				return nil
			}
			err := xorm.Update(tx, setting, updateMap, xorm.Where("convId = ? AND userId = ?", convId, in.CommonReq.UserId))
			if err != nil {
				l.Errorf("SetConvSetting: %v", err)
			}
			return err
		}
	}, func(tx *gorm.DB) error {
		notice := &noticemodel.Notice{
			ConvId: pb.HiddenConvId(convId),
			Options: noticemodel.NoticeOption{
				StorageForClient: false,
				UpdateConvNotice: false,
			},
			ContentType: pb.NoticeContentType_SyncConvSetting,
			Content: utils.AnyToBytes(pb.NoticeContent_SyncConvSetting{
				ConvIds: []string{convId},
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
			ConvIds:   []string{convId},
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
			ConvId:    pb.HiddenConvId(convId),
			DeviceId:  nil,
		})
		if err != nil {
			l.Errorf("GetUserNoticeData: %v", err)
		}
		return err
	})
	return &pb.UpdateConvSettingResp{}, nil
}
