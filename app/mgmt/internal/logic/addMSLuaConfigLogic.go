package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"
	"time"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddMSLuaConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddMSLuaConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddMSLuaConfigLogic {
	return &AddMSLuaConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddMSLuaConfigLogic) AddMSLuaConfig(in *pb.AddMSLuaConfigReq) (*pb.AddMSLuaConfigResp, error) {
	model := &mgmtmodel.LuaConfig{
		Id:         mgmtmodel.GetId(l.svcCtx.Mysql(), &mgmtmodel.LuaConfig{}, 1000),
		Name:       in.LuaConfig.Name,
		Desc:       in.LuaConfig.Desc,
		Code:       in.LuaConfig.Code,
		Type:       in.LuaConfig.Type,
		Enable:     in.LuaConfig.Enable,
		CreateTime: time.Now().UnixMilli(),
	}
	err := l.svcCtx.Mysql().Model(model).Create(model).Error
	if err != nil {
		l.Errorf("添加失败: %v", err)
		return &pb.AddMSLuaConfigResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	return &pb.AddMSLuaConfigResp{}, nil
}
