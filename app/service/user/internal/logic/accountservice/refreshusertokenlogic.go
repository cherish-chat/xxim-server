package accountservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/user/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefreshUserTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefreshUserTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefreshUserTokenLogic {
	return &RefreshUserTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RefreshUserToken 刷新用户token
func (l *RefreshUserTokenLogic) RefreshUserToken(in *peerpb.RefreshUserTokenReq) (*peerpb.RefreshUserTokenResp, error) {
	// todo: add your logic here and delete this line

	return &peerpb.RefreshUserTokenResp{}, nil
}
