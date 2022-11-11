package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/group/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type BanAllGroupMemberLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBanAllGroupMemberLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BanAllGroupMemberLogic {
	return &BanAllGroupMemberLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// BanAllGroupMember 禁言全部群成员
func (l *BanAllGroupMemberLogic) BanAllGroupMember(in *pb.BanAllGroupMemberReq) (*pb.BanAllGroupMemberResp, error) {
	// todo: add your logic here and delete this line

	return &pb.BanAllGroupMemberResp{}, nil
}
