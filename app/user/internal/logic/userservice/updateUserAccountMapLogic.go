package userservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

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
func (l *UpdateUserAccountMapLogic) UpdateUserAccountMap(in *pb.UpdateUserAccountMapReq) (*pb.UpdateUserAccountMapResp, error) {
	// todo: add your logic here and delete this line

	return &pb.UpdateUserAccountMapResp{}, nil
}
