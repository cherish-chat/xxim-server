package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/xx/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type KickGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewKickGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *KickGroupMemberLogic {
	return &KickGroupMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// KickGroupMember 踢出群组
func (l *KickGroupMemberLogic) KickGroupMember(in *pb.KickGroupMemberReq) (*pb.KickGroupMemberResp, error) {
	// todo: add your logic here and delete this line

	return &pb.KickGroupMemberResp{}, nil
}
