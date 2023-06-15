package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserOfflineCallbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserOfflineCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserOfflineCallbackLogic {
	return &UserOfflineCallbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UserOfflineCallback 用户下线回调
func (l *UserOfflineCallbackLogic) UserOfflineCallback(in *pb.UserOfflineCallbackReq) (*pb.UserOfflineCallbackResp, error) {
	// todo: add your logic here and delete this line

	return &pb.UserOfflineCallbackResp{}, nil
}
