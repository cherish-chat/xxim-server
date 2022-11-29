package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/im/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/xconf"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
)

type GetAppSystemConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAppSystemConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAppSystemConfigLogic {
	return &GetAppSystemConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAppSystemConfigLogic) GetAppSystemConfig(in *pb.GetAppSystemConfigReq) (*pb.GetAppSystemConfigResp, error) {
	var appSystemConfigs []*xconf.SystemConfig
	err := l.svcCtx.Mysql().Model(&xconf.SystemConfig{}).Where("`key` LIKE ?", "app.%").Find(&appSystemConfigs).Error
	if err != nil {
		l.Errorf("获取系统配置失败: %v", err)
		return &pb.GetAppSystemConfigResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	configs := make(map[string]string)
	for _, config := range appSystemConfigs {
		key := strings.TrimPrefix(config.Key, "app.")
		configs[key] = config.Value
	}
	return &pb.GetAppSystemConfigResp{
		CommonResp: pb.NewSuccessResp(),
		Configs:    configs,
	}, nil
}
