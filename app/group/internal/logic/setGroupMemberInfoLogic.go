package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetGroupMemberInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSetGroupMemberInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetGroupMemberInfoLogic {
	return &SetGroupMemberInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SetGroupMemberInfo 设置群成员信息
func (l *SetGroupMemberInfoLogic) SetGroupMemberInfo(in *pb.SetGroupMemberInfoReq) (*pb.SetGroupMemberInfoResp, error) {
	// todo: add your logic here and delete this line

	return &pb.SetGroupMemberInfoResp{}, nil
}
