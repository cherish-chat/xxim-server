package userservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserAfterOnlineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserAfterOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAfterOnlineLogic {
	return &UserAfterOnlineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UserAfterOnline 用户上线回调
func (l *UserAfterOnlineLogic) UserAfterOnline(in *pb.UserAfterOnlineReq) (*pb.UserAfterOnlineResp, error) {
	// todo: add your logic here and delete this line

	return &pb.UserAfterOnlineResp{}, nil
}
