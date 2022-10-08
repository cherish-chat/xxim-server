package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/xx/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type InviteGroupLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewInviteGroupLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InviteGroupLogic {
	return &InviteGroupLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// InviteGroup 邀请加入群组
func (l *InviteGroupLogic) InviteGroup(in *pb.InviteGroupReq) (*pb.InviteGroupResp, error) {
	// todo: add your logic here and delete this line

	return &pb.InviteGroupResp{}, nil
}
