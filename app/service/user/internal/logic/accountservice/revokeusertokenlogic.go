package accountservicelogic

import (
	"context"
	"github.com/cherish-chat/xxim-server/app/service/user/usermodel"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/user/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RevokeUserTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRevokeUserTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RevokeUserTokenLogic {
	return &RevokeUserTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RevokeUserToken 注销用户token
func (l *RevokeUserTokenLogic) RevokeUserToken(in *peerpb.RevokeUserTokenReq) (*peerpb.RevokeUserTokenResp, error) {
	err := l.svcCtx.Jwt.RevokeToken(context.Background(), in.Header.UserId, usermodel.GetJwtUniqueKey(in.Header))
	if err != nil {
		l.Errorf("revoke token error: %v", err)
		return &peerpb.RevokeUserTokenResp{}, err
	}
	return &peerpb.RevokeUserTokenResp{}, nil
}
