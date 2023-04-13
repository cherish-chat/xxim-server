package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMSLuaConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateMSLuaConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMSLuaConfigLogic {
	return &UpdateMSLuaConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateMSLuaConfigLogic) UpdateMSLuaConfig(in *pb.UpdateMSLuaConfigReq) (*pb.UpdateMSLuaConfigResp, error) {
	// 查询原模型
	model := &mgmtmodel.LuaConfig{}
	err := l.svcCtx.Mysql().Model(model).Where("id = ?", in.LuaConfig.Id).First(model).Error
	if err != nil {
		l.Errorf("查询失败: %v", err)
		return &pb.UpdateMSLuaConfigResp{CommonResp: pb.NewRetryErrorResp()}, err
	}
	updateMap := make(map[string]interface{})
	if in.LuaConfig.Name != "" {
		updateMap["name"] = in.LuaConfig.Name
	}
	if in.LuaConfig.Desc != "" {
		updateMap["desc"] = in.LuaConfig.Desc
	}
	if in.LuaConfig.Code != "" {
		updateMap["code"] = in.LuaConfig.Code
	}
	updateMap["type"] = in.LuaConfig.Type
	updateMap["enable"] = in.LuaConfig.Enable
	if len(updateMap) > 0 {
		err = l.svcCtx.Mysql().Model(model).Where("id = ?", in.LuaConfig.Id).Updates(updateMap).Error
		if err != nil {
			l.Errorf("更新失败: %v", err)
			return &pb.UpdateMSLuaConfigResp{CommonResp: pb.NewRetryErrorResp()}, err
		}
	}
	return &pb.UpdateMSLuaConfigResp{}, nil
}
