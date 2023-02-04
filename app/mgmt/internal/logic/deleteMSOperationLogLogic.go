package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMSOperationLogLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMSOperationLogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMSOperationLogLogic {
	return &DeleteMSOperationLogLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteMSOperationLogLogic) DeleteMSOperationLog(in *pb.DeleteMSOperationLogReq) (*pb.DeleteMSOperationLogResp, error) {
	err := l.svcCtx.Mysql().Model(&mgmtmodel.OperationLog{}).Where("id in (?)", in.Ids).Delete(&mgmtmodel.OperationLog{}).Error
	if err != nil {
		l.Errorf("delete error: %v", err)
	}
	return &pb.DeleteMSOperationLogResp{}, err
}
