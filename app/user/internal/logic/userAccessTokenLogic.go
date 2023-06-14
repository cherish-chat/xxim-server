package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserAccessTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserAccessTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAccessTokenLogic {
	return &UserAccessTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UserAccessToken 用户登录
func (l *UserAccessTokenLogic) UserAccessToken(in *pb.UserAccessTokenReq) (*pb.UserAccessTokenResp, error) {
	// todo: add your logic here and delete this line

	return &pb.UserAccessTokenResp{}, nil
}
