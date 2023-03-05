package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"
	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateAppMgmtConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateAppMgmtConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateAppMgmtConfigLogic {
	return &UpdateAppMgmtConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateAppMgmtConfigLogic) UpdateAppMgmtConfig(in *pb.UpdateAppMgmtConfigReq) (*pb.UpdateAppMgmtConfigResp, error) {
	err := l.svcCtx.ConfigMgr.Flush(l.ctx, "")
	if err != nil {
		l.Errorf("flush config failed, err: %v", err)
		return &pb.UpdateAppMgmtConfigResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	err = l.svcCtx.ConfigMgr.Flush(l.ctx, in.UserId)
	if err != nil {
		l.Errorf("flush config failed, err: %v", err)
		return &pb.UpdateAppMgmtConfigResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	now := time.Now().UnixMilli()
	for _, config := range in.AppMgmtConfigs {
		model := &appmgmtmodel.Config{
			Group:          config.Group,
			K:              config.K,
			V:              config.V,
			Type:           config.Type,
			Name:           config.Name,
			ScopePlatforms: config.ScopePlatforms,
			UserId:         in.UserId,
		}
		err := l.upsert(model, map[string]any{
			"v":          config.V,
			"updateTime": now,
		})
		if err != nil {
			l.Errorf("update config failed, err: %v", err)
			return &pb.UpdateAppMgmtConfigResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	go l.svcCtx.ConfigMgr.Flush(l.ctx, in.UserId)
	return &pb.UpdateAppMgmtConfigResp{}, nil
}

func (l *UpdateAppMgmtConfigLogic) upsert(model *appmgmtmodel.Config, m map[string]any) error {
	// 使用 k 和 userId 作为唯一索引
	// 如果存在则更新，否则插入
	tx := l.svcCtx.Mysql().Model(model).Where("k = ? AND userId = ?", model.K, model.UserId).Updates(m)
	err := tx.Error
	if err != nil {
		return err
	}
	if tx.RowsAffected == 0 {
		conf := &appmgmtmodel.Config{}
		err := l.svcCtx.Mysql().Model(conf).Where("k = ? AND userId = ?", model.K, "").First(conf).Error
		if err != nil {
			conf = model
		}
		return l.svcCtx.Mysql().Model(&appmgmtmodel.Config{}).Create(&appmgmtmodel.Config{
			Group:          conf.Group,
			K:              conf.K,
			V:              model.V,
			Type:           conf.Type,
			Name:           conf.Name,
			ScopePlatforms: conf.ScopePlatforms,
			Options:        conf.Options,
			UserId:         model.UserId,
			UpdateTime:     time.Now().UnixMilli(),
		}).Error
	}
	return nil
}
