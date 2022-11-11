package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnbanGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnbanGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnbanGroupMemberLogic {
	return &UnbanGroupMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UnbanGroupMember 解除禁言群成员
func (l *UnbanGroupMemberLogic) UnbanGroupMember(in *pb.UnbanGroupMemberReq) (*pb.UnbanGroupMemberResp, error) {
	// todo: add your logic here and delete this line

	return &pb.UnbanGroupMemberResp{}, nil
}
