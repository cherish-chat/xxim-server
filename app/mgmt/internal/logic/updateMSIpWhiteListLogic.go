package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMSIpWhiteListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateMSIpWhiteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMSIpWhiteListLogic {
	return &UpdateMSIpWhiteListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateMSIpWhiteListLogic) UpdateMSIpWhiteList(in *pb.UpdateMSIpWhiteListReq) (*pb.UpdateMSIpWhiteListResp, error) {
	// todo: add your logic here and delete this line

	return &pb.UpdateMSIpWhiteListResp{}, nil
}
