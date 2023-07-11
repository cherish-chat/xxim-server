package callbackservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/user/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserAfterOnlineLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserAfterOnlineLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAfterOnlineLogic {
	return &UserAfterOnlineLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UserAfterOnlineLogic) UserAfterOnline(in *peerpb.UserAfterOnlineReq) (*peerpb.UserAfterOnlineResp, error) {
	// todo: add your logic here and delete this line

	return &peerpb.UserAfterOnlineResp{}, nil
}
