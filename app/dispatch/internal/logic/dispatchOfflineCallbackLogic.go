package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/dispatch/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DispatchOfflineCallbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDispatchOfflineCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DispatchOfflineCallbackLogic {
	return &DispatchOfflineCallbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// DispatchOfflineCallback 下线回调
func (l *DispatchOfflineCallbackLogic) DispatchOfflineCallback(in *pb.DispatchOfflineCallbackReq) (*pb.DispatchOfflineCallbackResp, error) {
	// todo: add your logic here and delete this line

	return &pb.DispatchOfflineCallbackResp{}, nil
}
