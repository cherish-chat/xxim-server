package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type KeepAliveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewKeepAliveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *KeepAliveLogic {
	return &KeepAliveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *KeepAliveLogic) KeepAlive(in *pb.KeepAliveReq) (*pb.KeepAliveResp, error) {
	_, err := NewFlushUsersSubConvLogic(l.ctx, l.svcCtx).FlushUsersSubConv(&pb.FlushUsersSubConvReq{UserIds: []string{in.CommonReq.UserId}})
	if err != nil {
		return &pb.KeepAliveResp{}, err
	}
	return &pb.KeepAliveResp{}, nil
}
