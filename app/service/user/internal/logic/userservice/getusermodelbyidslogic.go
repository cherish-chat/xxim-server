package userservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/user/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserModelByIdsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserModelByIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserModelByIdsLogic {
	return &GetUserModelByIdsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetUserModelByIds 批量获取用户模型
func (l *GetUserModelByIdsLogic) GetUserModelByIds(in *peerpb.GetUserModelByIdsReq) (*peerpb.GetUserModelByIdsResp, error) {
	// todo: add your logic here and delete this line

	return &peerpb.GetUserModelByIdsResp{}, nil
}
