package logic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/xx/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetFriendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetFriendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetFriendLogic {
	return &GetFriendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetFriend 获取好友
func (l *GetFriendLogic) GetFriend(in *pb.GetFriendReq) (*pb.GetFriendResp, error) {
	// todo: add your logic here and delete this line

	return &pb.GetFriendResp{}, nil
}
