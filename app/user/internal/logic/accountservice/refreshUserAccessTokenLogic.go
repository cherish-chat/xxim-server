package accountservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshUserAccessTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefreshUserAccessTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshUserAccessTokenLogic {
	return &RefreshUserAccessTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// deprecated
// RefreshUserAccessToken 刷新用户token
func (l *RefreshUserAccessTokenLogic) RefreshUserAccessToken(in *pb.RefreshUserAccessTokenReq) (*pb.RefreshUserAccessTokenResp, error) {
	return &pb.RefreshUserAccessTokenResp{}, nil
}
