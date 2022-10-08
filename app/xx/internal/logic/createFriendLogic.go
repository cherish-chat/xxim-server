package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/xx/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateFriendLogic {
	return &CreateFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CreateFriend 添加好友
func (l *CreateFriendLogic) CreateFriend(in *pb.CreateFriendReq) (*pb.CreateFriendResp, error) {
	// todo: add your logic here and delete this line

	return &pb.CreateFriendResp{}, nil
}
