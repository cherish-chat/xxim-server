package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupMemberListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupMemberListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMemberListLogic {
	return &GetGroupMemberListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetGroupMemberList 获取群成员列表
func (l *GetGroupMemberListLogic) GetGroupMemberList(in *pb.GetGroupMemberListReq) (*pb.GetGroupMemberListResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetGroupMemberListResp{}, nil
}
