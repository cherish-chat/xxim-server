package userservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-proto/peerpb"
	"github.com/cherish-chat/xxim-server/app/service/user/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserModelByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserModelByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserModelByIdLogic {
	return &GetUserModelByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetUserModelById 获取用户模型
func (l *GetUserModelByIdLogic) GetUserModelById(in *peerpb.GetUserModelByIdReq) (*peerpb.GetUserModelByIdResp, error) {
	// todo: add your logic here and delete this line

	return &peerpb.GetUserModelByIdResp{}, nil
}
