package callbackservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/user/internal/svc"

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

func (l *UserBeforeRequestLogic) UserBeforeRequest(in *peerpb.UserBeforeRequestReq) (*peerpb.UserBeforeRequestResp, error) {
	return &peerpb.UserBeforeRequestResp{
		Header: peerpb.NewOkHeader(),
	}, nil
}
