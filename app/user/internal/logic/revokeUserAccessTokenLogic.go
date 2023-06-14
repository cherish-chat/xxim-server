package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RevokeUserAccessTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRevokeUserAccessTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RevokeUserAccessTokenLogic {
	return &RevokeUserAccessTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RevokeUserAccessToken 注销用户token
func (l *RevokeUserAccessTokenLogic) RevokeUserAccessToken(in *pb.RevokeUserAccessTokenReq) (*pb.RevokeUserAccessTokenResp, error) {
	// todo: add your logic here and delete this line

	return &pb.RevokeUserAccessTokenResp{}, nil
}
