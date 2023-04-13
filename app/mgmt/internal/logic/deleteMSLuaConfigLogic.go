package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMSLuaConfigLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMSLuaConfigLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMSLuaConfigLogic {
	return &DeleteMSLuaConfigLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteMSLuaConfigLogic) DeleteMSLuaConfig(in *pb.DeleteMSLuaConfigReq) (*pb.DeleteMSLuaConfigResp, error) {
	err := l.svcCtx.Mysql().Model(&mgmtmodel.LuaConfig{}).Where("id in (?)", in.Ids).Delete(&mgmtmodel.LuaConfig{}).Error
	if err != nil {
		l.Errorf("delete error: %v", err)
	}
	return &pb.DeleteMSLuaConfigResp{}, err
}
