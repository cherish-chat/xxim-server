package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/xx/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ClearGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewClearGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ClearGroupMemberLogic {
	return &ClearGroupMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ClearGroupMember 清空群组成员
func (l *ClearGroupMemberLogic) ClearGroupMember(in *pb.ClearGroupMemberReq) (*pb.ClearGroupMemberResp, error) {
	// todo: add your logic here and delete this line

	return &pb.ClearGroupMemberResp{}, nil
}
