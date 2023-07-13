package accountservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/user/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserDestroyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserDestroyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserDestroyLogic {
	return &UserDestroyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UserDestroy 用户注销
func (l *UserDestroyLogic) UserDestroy(in *peerpb.UserDestroyReq) (*peerpb.UserDestroyResp, error) {
	// todo: add your logic here and delete this line

	return &peerpb.UserDestroyResp{}, nil
}
