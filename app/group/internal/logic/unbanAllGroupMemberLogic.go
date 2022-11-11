package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnbanAllGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnbanAllGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnbanAllGroupMemberLogic {
	return &UnbanAllGroupMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UnbanAllGroupMember 解除禁言全部群成员
func (l *UnbanAllGroupMemberLogic) UnbanAllGroupMember(in *pb.UnbanAllGroupMemberReq) (*pb.UnbanAllGroupMemberResp, error) {
	// todo: add your logic here and delete this line

	return &pb.UnbanAllGroupMemberResp{}, nil
}
