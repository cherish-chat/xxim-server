package callbackservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/user/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserBeforeConnectLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserBeforeConnectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserBeforeConnectLogic {
	return &UserBeforeConnectLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserBeforeConnectLogic) UserBeforeConnect(in *peerpb.UserBeforeConnectReq) (*peerpb.UserBeforeConnectResp, error) {
	// todo: add your logic here and delete this line

	return &peerpb.UserBeforeConnectResp{}, nil
}
