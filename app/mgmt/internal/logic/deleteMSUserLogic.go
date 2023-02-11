package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMSUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMSUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMSUserLogic {
	return &DeleteMSUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteMSUserLogic) DeleteMSUser(in *pb.DeleteMSUserReq) (*pb.DeleteMSUserResp, error) {
	err := l.svcCtx.Mysql().Model(&mgmtmodel.User{}).Where("id in (?)", in.Ids).Delete(&mgmtmodel.User{}).Error
	if err != nil {
		l.Errorf("delete user error: %v", err)
	}
	return &pb.DeleteMSUserResp{}, err
}
