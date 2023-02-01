package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/mgmt/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMSApiPathDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMSApiPathDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMSApiPathDetailLogic {
	return &GetMSApiPathDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetMSApiPathDetailLogic) GetMSApiPathDetail(in *pb.GetMSApiPathDetailReq) (*pb.GetMSApiPathDetailResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetMSApiPathDetailResp{}, nil
}
