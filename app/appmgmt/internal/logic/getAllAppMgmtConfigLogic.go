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
	models, err := l.svcCtx.ConfigMgr.GetAll(l.ctx, "")
	if err != nil {
		l.Errorf("get all app mgmt config error: %v", err)
		return &pb.GetAllAppMgmtConfigResp{
			CommonResp: pb.NewRetryErrorResp(),
		}, err
	}
	var resp []*pb.AppMgmtConfig
	var respKvMap = make(map[string]*pb.AppMgmtConfig)
	for _, model := range models {
		toPB := model.ToPB()
		resp = append(resp, toPB)
		respKvMap[model.K] = toPB
	}
	if in.UserId != "" {
		// 单独查询用户配置
		configs, err := l.svcCtx.ConfigMgr.GetAll(l.ctx, in.UserId)
		if err != nil {
			l.Errorf("get all app mgmt config error: %v", err)
			return &pb.GetAllAppMgmtConfigResp{
				CommonResp: pb.NewRetryErrorResp(),
			}, err
		}
		resp = []*pb.AppMgmtConfig{}
		for _, config := range configs {
			if _, ok := respKvMap[config.K]; !ok {
				respKvMap[config.K] = config.ToPB()
			}
		}
		for _, v := range respKvMap {
			resp = append(resp, v)
		}
	}
	return &pb.GetAllAppMgmtConfigResp{
		AppMgmtConfigs: resp,
	}, nil
}
