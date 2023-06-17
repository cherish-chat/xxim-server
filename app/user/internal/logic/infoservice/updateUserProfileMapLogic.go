package infoservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/user/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserProfileMapLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserProfileMapLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserProfileMapLogic {
	return &UpdateUserProfileMapLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateUserProfileMap 更新用户个人信息
func (l *UpdateUserProfileMapLogic) UpdateUserProfileMap(in *pb.UpdateUserProfileMapReq) (*pb.UpdateUserProfileMapResp, error) {
	// todo: add your logic here and delete this line

	return &pb.UpdateUserProfileMapResp{}, nil
}
