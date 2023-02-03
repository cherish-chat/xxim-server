package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMSMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMSMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMSMenuLogic {
	return &DeleteMSMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteMSMenuLogic) DeleteMSMenu(in *pb.DeleteMSMenuReq) (*pb.DeleteMSMenuResp, error) {
	err := l.svcCtx.Mysql().Model(&mgmtmodel.Menu{}).Where("id in (?)", in.Ids).Delete(&mgmtmodel.Menu{}).Error
	if err != nil {
		l.Errorf("delete error: %v", err)
	}
	return &pb.DeleteMSMenuResp{}, err
}
