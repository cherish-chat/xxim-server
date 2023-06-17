package callbackservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserAfterOfflineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserAfterOfflineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAfterOfflineLogic {
	return &UserAfterOfflineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UserAfterOffline 用户下线回调
func (l *UserAfterOfflineLogic) UserAfterOffline(in *pb.UserAfterOfflineReq) (*pb.UserAfterOfflineResp, error) {
	// todo: add your logic here and delete this line

	return &pb.UserAfterOfflineResp{}, nil
}
