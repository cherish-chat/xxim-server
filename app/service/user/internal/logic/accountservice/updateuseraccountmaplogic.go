package accountservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/user/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserAccountMapLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserAccountMapLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserAccountMapLogic {
	return &UpdateUserAccountMapLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateUserAccountMap 更新用户账号信息
func (l *UpdateUserAccountMapLogic) UpdateUserAccountMap(in *peerpb.UpdateUserAccountMapReq) (*peerpb.UpdateUserAccountMapResp, error) {
	// todo: add your logic here and delete this line

	return &peerpb.UpdateUserAccountMapResp{}, nil
}
