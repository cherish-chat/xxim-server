package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetGroupMemberRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetGroupMemberRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetGroupMemberRoleLogic {
	return &SetGroupMemberRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SetGroupMemberRole 设置群成员角色
func (l *SetGroupMemberRoleLogic) SetGroupMemberRole(in *pb.SetGroupMemberRoleReq) (*pb.SetGroupMemberRoleResp, error) {
	// todo: add your logic here and delete this line

	return &pb.SetGroupMemberRoleResp{}, nil
}
