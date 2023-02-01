package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddMSIpWhiteListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddMSIpWhiteListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddMSIpWhiteListLogic {
	return &AddMSIpWhiteListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddMSIpWhiteListLogic) AddMSIpWhiteList(in *pb.AddMSIpWhiteListReq) (*pb.AddMSIpWhiteListResp, error) {
	// todo: add your logic here and delete this line

	return &pb.AddMSIpWhiteListResp{}, nil
}
