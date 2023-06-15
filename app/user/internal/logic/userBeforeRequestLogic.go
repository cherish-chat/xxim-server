package logic

import (
	"context"
	"github.com/cherish-chat/xxim-server/common/i18n"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserBeforeRequestLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserBeforeRequestLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserBeforeRequestLogic {
	return &UserBeforeRequestLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UserBeforeRequest 用户请求前的回调
func (l *UserBeforeRequestLogic) UserBeforeRequest(in *pb.UserBeforeRequestReq) (*pb.UserBeforeRequestResp, error) {
	// todo: add your logic here and delete this line
	if in.Header.UserToken != "" {
		tokenObject, err := l.svcCtx.Jwt.VerifyToken(l.ctx, in.Header.UserToken, in.Header.GetJwtUniqueKey())
		if err != nil {
			l.Errorf("verifyToken error: %s", err.Error())
			return nil, err
		}
		l.Debugf("tokenObject: %+v", tokenObject)
		return &pb.UserBeforeRequestResp{
			Header: i18n.NewOkHeader(),
		}, nil
	}
	return &pb.UserBeforeRequestResp{
		Header: i18n.NewOkHeader(),
	}, nil
}
