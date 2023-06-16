package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/i18n"
	"time"

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

// RefreshUserAccessToken 刷新用户token
func (l *RefreshUserAccessTokenLogic) RefreshUserAccessToken(in *pb.RefreshUserAccessTokenReq) (*pb.RefreshUserAccessTokenResp, error) {
	tokenObject, verifyTokenErr := l.svcCtx.Jwt.VerifyToken(l.ctx, in.Header.UserToken, in.Header.GetJwtUniqueKey())
	if verifyTokenErr != nil {
		l.Errorf("verifyTokenErr: %v", verifyTokenErr)
		return &pb.RefreshUserAccessTokenResp{
			Header: i18n.NewAuthError(pb.AuthErrorTypeExpired, ""),
		}, nil
	}
	tokenObject.ExpiredAt = time.Now().Add(time.Hour * time.Duration(l.svcCtx.Config.Account.JwtConfig.ExpireHour)).UnixMilli()
	tokenObject.AliveTime = time.Now().UnixMilli()
	err := l.svcCtx.Jwt.SetToken(l.ctx, tokenObject)
	if err != nil {
		l.Errorf("setTokenErr: %v", err)
		return nil, err
	}
	return &pb.RefreshUserAccessTokenResp{}, nil
}
