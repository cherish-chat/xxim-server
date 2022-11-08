package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/relation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type AreFriendsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAreFriendsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AreFriendsLogic {
	return &AreFriendsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *AreFriendsLogic) AreFriends(in *pb.AreFriendsReq) (*pb.AreFriendsResp, error) {
	// todo: add your logic here and delete this line

	return &pb.AreFriendsResp{}, nil
}
