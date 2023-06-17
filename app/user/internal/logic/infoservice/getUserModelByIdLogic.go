package infoservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

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
func (l *GetUserModelByIdLogic) GetUserModelById(in *pb.GetUserModelByIdReq) (*pb.GetUserModelByIdResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetUserModelByIdResp{}, nil
}
