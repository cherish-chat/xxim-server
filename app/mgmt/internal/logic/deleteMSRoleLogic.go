package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMSRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMSRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMSRoleLogic {
	return &DeleteMSRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteMSRoleLogic) DeleteMSRole(in *pb.DeleteMSRoleReq) (*pb.DeleteMSRoleResp, error) {
	err := l.svcCtx.Mysql().Model(&mgmtmodel.Role{}).Where("id in (?)", in.Ids).Delete(&mgmtmodel.Role{}).Error
	if err != nil {
		l.Errorf("delete error: %v", err)
	}
	return &pb.DeleteMSRoleResp{}, err
}
