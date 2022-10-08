package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/xx/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateFriendLogic {
	return &UpdateFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateFriend 更新好友
func (l *UpdateFriendLogic) UpdateFriend(in *pb.UpdateFriendReq) (*pb.UpdateFriendResp, error) {
	// todo: add your logic here and delete this line

	return &pb.UpdateFriendResp{}, nil
}
