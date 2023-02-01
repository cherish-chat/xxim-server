package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddMSApiPathLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddMSApiPathLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddMSApiPathLogic {
	return &AddMSApiPathLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddMSApiPathLogic) AddMSApiPath(in *pb.AddMSApiPathReq) (*pb.AddMSApiPathResp, error) {
	// todo: add your logic here and delete this line

	return &pb.AddMSApiPathResp{}, nil
}
