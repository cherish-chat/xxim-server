package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

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
	err := l.svcCtx.ConfigMgr.Flush(l.ctx)
	if err != nil {
		l.Errorf("flush config failed, err: %v", err)
		return &pb.UpdateAppMgmtConfigResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	for _, config := range in.AppMgmtConfigs {
		model := &appmgmtmodel.Config{
			Group:          config.Group,
			K:              config.K,
			V:              config.V,
			Type:           config.Type,
			Name:           config.Name,
			ScopePlatforms: config.ScopePlatforms,
		}
		err := l.svcCtx.Mysql().Model(model).Where("k = ?", config.K).Updates(map[string]any{
			"group": config.Group,
			"v":     config.V,
		}).Error
		if err != nil {
			l.Errorf("update config failed, err: %v", err)
			return &pb.UpdateAppMgmtConfigResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	err = l.svcCtx.ConfigMgr.Flush(l.ctx)
	if err != nil {
		l.Errorf("flush config failed, err: %v", err)
		return &pb.UpdateAppMgmtConfigResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.UpdateAppMgmtConfigResp{}, nil
}
