package userservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/user/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserCountMapLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserCountMapLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserCountMapLogic {
	return &UpdateUserCountMapLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateUserCountMap 更新用户计数信息
func (l *UpdateUserCountMapLogic) UpdateUserCountMap(in *peerpb.UpdateUserCountMapReq) (*peerpb.UpdateUserCountMapResp, error) {
	// todo: add your logic here and delete this line

	return &peerpb.UpdateUserCountMapResp{}, nil
}
