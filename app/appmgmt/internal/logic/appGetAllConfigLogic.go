package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AppGetAllConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAppGetAllConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AppGetAllConfigLogic {
	return &AppGetAllConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AppGetAllConfigLogic) AppGetAllConfig(in *pb.AppGetAllConfigReq) (*pb.AppGetAllConfigResp, error) {
	configs, err := l.svcCtx.ConfigMgr.GetAll(l.ctx, "")
	if err != nil {
		l.Errorf("get all app config error: %v", err)
		return &pb.AppGetAllConfigResp{
			CommonResp: pb.NewRetryErrorResp(),
		}, err
	}
	var configMap = make(map[string]string)
	for _, config := range configs {
		configMap[config.K] = config.V
	}
	configs, err = l.svcCtx.ConfigMgr.GetAll(l.ctx, in.CommonReq.UserId)
	if err != nil {
		l.Errorf("get all app mgmt config error: %v", err)
		return &pb.AppGetAllConfigResp{ConfigMap: configMap}, nil
	}
	for _, config := range configs {
		configMap[config.K] = config.V
	}
	return &pb.AppGetAllConfigResp{ConfigMap: configMap}, nil
}
