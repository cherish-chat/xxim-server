package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserOnlineCallbackLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserOnlineCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserOnlineCallbackLogic {
	return &UserOnlineCallbackLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UserOnlineCallback 用户上线回调
func (l *UserOnlineCallbackLogic) UserOnlineCallback(in *pb.UserOnlineCallbackReq) (*pb.UserOnlineCallbackResp, error) {
	// todo: add your logic here and delete this line

	return &pb.UserOnlineCallbackResp{}, nil
}
