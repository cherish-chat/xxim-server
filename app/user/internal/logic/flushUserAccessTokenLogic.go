package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FlushUserAccessTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFlushUserAccessTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FlushUserAccessTokenLogic {
	return &FlushUserAccessTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// FlushUserAccessToken 刷新用户token
func (l *FlushUserAccessTokenLogic) FlushUserAccessToken(in *pb.FlushUserAccessTokenReq) (*pb.FlushUserAccessTokenResp, error) {
	// todo: add your logic here and delete this line

	return &pb.FlushUserAccessTokenResp{}, nil
}
