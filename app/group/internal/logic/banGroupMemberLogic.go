package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type BanGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBanGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BanGroupMemberLogic {
	return &BanGroupMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BanGroupMember 禁言群成员
func (l *BanGroupMemberLogic) BanGroupMember(in *pb.BanGroupMemberReq) (*pb.BanGroupMemberResp, error) {
	// todo: add your logic here and delete this line

	return &pb.BanGroupMemberResp{}, nil
}
