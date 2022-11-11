package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGroupMemberInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGroupMemberInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGroupMemberInfoLogic {
	return &GetGroupMemberInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetGroupMemberInfo 获取群成员信息
func (l *GetGroupMemberInfoLogic) GetGroupMemberInfo(in *pb.GetGroupMemberInfoReq) (*pb.GetGroupMemberInfoResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetGroupMemberInfoResp{}, nil
}
