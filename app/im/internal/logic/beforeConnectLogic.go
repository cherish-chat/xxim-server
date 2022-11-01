package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/im/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type BeforeConnectLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBeforeConnectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BeforeConnectLogic {
	return &BeforeConnectLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BeforeConnectLogic) BeforeConnect(in *pb.BeforeConnectReq) (*pb.BeforeConnectResp, error) {
	// todo: add your logic here and delete this line

	return &pb.BeforeConnectResp{}, nil
}
