package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/msg/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"
	"github.com/cherish-chat/xxim-server/common/xtrace"

	"github.com/zeromicro/go-zero/core/logx"
)

type AfterConnectLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAfterConnectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AfterConnectLogic {
	return &AfterConnectLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// conn hook
func (l *AfterConnectLogic) AfterConnect(in *pb.AfterConnectReq) (*pb.CommonResp, error) {
	var err error
	xtrace.RunWithTrace(xtrace.TraceIdFromContext(l.ctx), "msgService/AfterConnect/FlushUsersSubConv", func(ctx context.Context) {
		_, err = NewFlushUsersSubConvLogic(ctx, l.svcCtx).FlushUsersSubConv(&pb.FlushUsersSubConvReq{UserIds: []string{in.ConnParam.UserId}})
	}, nil)
	if err != nil {
		return &pb.CommonResp{}, err
	}
	return &pb.CommonResp{}, nil
}
