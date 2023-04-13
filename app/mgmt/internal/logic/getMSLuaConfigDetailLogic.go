package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMSLuaConfigDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMSLuaConfigDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMSLuaConfigDetailLogic {
	return &GetMSLuaConfigDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMSLuaConfigDetailLogic) GetMSLuaConfigDetail(in *pb.GetMSLuaConfigReq) (*pb.GetMSLuaConfigResp, error) {
	// 查询原模型
	model := &mgmtmodel.LuaConfig{}
	err := l.svcCtx.Mysql().Model(model).Where("id = ?", in.Id).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.GetMSLuaConfigResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.GetMSLuaConfigResp{
		LuaConfig: model.ToPB(),
	}, nil
}
