package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/appmgmt/appmgmtmodel"

	"github.com/cherish-chat/xxim-server/app/appmgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteAppMgmtLinkLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteAppMgmtLinkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteAppMgmtLinkLogic {
	return &DeleteAppMgmtLinkLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteAppMgmtLinkLogic) DeleteAppMgmtLink(in *pb.DeleteAppMgmtLinkReq) (*pb.DeleteAppMgmtLinkResp, error) {
	model := &appmgmtmodel.Link{}
	err := l.svcCtx.Mysql().Model(model).Where("id in (?)", in.Ids).Delete(model).Error
	if err != nil {
		l.Errorf("delete error: %v", err)
		return &pb.DeleteAppMgmtLinkResp{
			CommonResp: pb.NewRetryErrorResp(),
		}, err
	}
	return &pb.DeleteAppMgmtLinkResp{}, nil
}
