package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDestroyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserDestroyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDestroyLogic {
	return &UserDestroyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UserDestroy 用户注销
func (l *UserDestroyLogic) UserDestroy(in *pb.UserDestroyReq) (*pb.UserDestroyResp, error) {
	// todo: add your logic here and delete this line

	return &pb.UserDestroyResp{}, nil
}
