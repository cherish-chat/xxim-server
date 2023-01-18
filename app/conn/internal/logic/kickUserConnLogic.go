package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/conn/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type KickUserConnLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewKickUserConnLogic(ctx context.Context, svcCtx *svc.ServiceContext) *KickUserConnLogic {
	return &KickUserConnLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *KickUserConnLogic) KickUserConn(in *pb.KickUserConnReq) (*pb.KickUserConnResp, error) {
	return &pb.KickUserConnResp{}, GetConnLogic().KickUserConn(l.ctx, in)
}
