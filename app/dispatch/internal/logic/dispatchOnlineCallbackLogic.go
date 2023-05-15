package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/dispatch/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DispatchOnlineCallbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDispatchOnlineCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DispatchOnlineCallbackLogic {
	return &DispatchOnlineCallbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DispatchOnlineCallback 上线回调
func (l *DispatchOnlineCallbackLogic) DispatchOnlineCallback(in *pb.DispatchOnlineCallbackReq) (*pb.DispatchOnlineCallbackResp, error) {
	// todo: add your logic here and delete this line

	return &pb.DispatchOnlineCallbackResp{}, nil
}
