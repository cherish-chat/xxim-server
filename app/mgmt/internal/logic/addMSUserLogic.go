package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddMSUserLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddMSUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddMSUserLogic {
	return &AddMSUserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AddMSUserLogic) AddMSUser(in *pb.AddMSUserReq) (*pb.AddMSUserResp, error) {
	// todo: add your logic here and delete this line

	return &pb.AddMSUserResp{}, nil
}
