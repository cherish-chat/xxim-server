package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/relation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RelationApplyFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRelationApplyFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RelationApplyFriendLogic {
	return &RelationApplyFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// RelationApplyFriend 申请好友
func (l *RelationApplyFriendLogic) RelationApplyFriend(in *pb.RelationApplyFriendReq) (*pb.RelationApplyFriendResp, error) {
	// todo: add your logic here and delete this line

	return &pb.RelationApplyFriendResp{}, nil
}
