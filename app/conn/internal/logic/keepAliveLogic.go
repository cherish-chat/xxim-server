package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/conn/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc"
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
	return &pb.KeepAliveResp{}, nil
}

var singleKeepAliveLogic *KeepAliveLogic

func GetKeepAliveLogic(svcCtx *svc.ServiceContext) *KeepAliveLogic {
	if singleKeepAliveLogic == nil {
		singleKeepAliveLogic = NewKeepAliveLogic(context.Background(), svcCtx)
	}
	return singleKeepAliveLogic
}

func (l *KeepAliveLogic) DoKeepAlive(ctx context.Context, req *pb.KeepAliveReq, opts ...grpc.CallOption) (*pb.KeepAliveResp, error) {
	return l.KeepAlive(req)
}
