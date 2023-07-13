package userservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/user/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserExtraMapLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserExtraMapLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserExtraMapLogic {
	return &UpdateUserExtraMapLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateUserExtraMap 更新用户扩展信息
func (l *UpdateUserExtraMapLogic) UpdateUserExtraMap(in *peerpb.UpdateUserExtraMapReq) (*peerpb.UpdateUserExtraMapResp, error) {
	// todo: add your logic here and delete this line

	return &peerpb.UpdateUserExtraMapResp{}, nil
}
