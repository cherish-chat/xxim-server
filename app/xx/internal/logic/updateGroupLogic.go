package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/xx/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateGroupLogic {
	return &UpdateGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateGroup 更新群组
func (l *UpdateGroupLogic) UpdateGroup(in *pb.UpdateGroupReq) (*pb.UpdateGroupResp, error) {
	// todo: add your logic here and delete this line

	return &pb.UpdateGroupResp{}, nil
}
