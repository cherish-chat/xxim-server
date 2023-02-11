package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAllAppMgmtConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAllAppMgmtConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAllAppMgmtConfigLogic {
	return &GetAllAppMgmtConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetAllAppMgmtConfigLogic) GetAllAppMgmtConfig(in *pb.GetAllAppMgmtConfigReq) (*pb.GetAllAppMgmtConfigResp, error) {
	models, err := l.svcCtx.ConfigMgr.GetAll(l.ctx)
	if err != nil {
		l.Errorf("get all app mgmt config error: %v", err)
		return &pb.GetAllAppMgmtConfigResp{
			CommonResp: pb.NewRetryErrorResp(),
		}, err
	}
	//var models []*appmgmtmodel.Config
	//l.svcCtx.Mysql().Model(&appmgmtmodel.Config{}).Find(&models)
	var resp []*pb.AppMgmtConfig
	for _, model := range models {
		resp = append(resp, model.ToPB())
	}
	return &pb.GetAllAppMgmtConfigResp{
		AppMgmtConfigs: resp,
	}, nil
}
