package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/mgmt/mgmtmodel"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteMSIpWhiteListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteMSIpWhiteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMSIpWhiteListLogic {
	return &DeleteMSIpWhiteListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteMSIpWhiteListLogic) DeleteMSIpWhiteList(in *pb.DeleteMSIpWhiteListReq) (*pb.DeleteMSIpWhiteListResp, error) {
	err := l.svcCtx.Mysql().Model(&mgmtmodel.MSIPWhitelist{}).Where("id in (?)", in.Ids).Delete(&mgmtmodel.MSIPWhitelist{}).Error
	if err != nil {
		l.Errorf("delete error: %v", err)
	}
	return &pb.DeleteMSIpWhiteListResp{}, err
}
