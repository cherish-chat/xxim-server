package callbackservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/user/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserAfterOfflineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserAfterOfflineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAfterOfflineLogic {
	return &UserAfterOfflineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserAfterOfflineLogic) UserAfterOffline(in *peerpb.UserAfterOfflineReq) (*peerpb.UserAfterOfflineResp, error) {
	// todo: add your logic here and delete this line

	return &peerpb.UserAfterOfflineResp{}, nil
}
