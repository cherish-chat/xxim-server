package friendservicelogic

import (
	"context"

	"github.com/cherish-chat/xxim-server/app/conversation/internal/svc"
	"github.com/cherish-chat/xxim-server/common/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type FriendApplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewFriendApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FriendApplyLogic {
	return &FriendApplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// FriendApply 添加好友
func (l *FriendApplyLogic) FriendApply(in *pb.FriendApplyReq) (*pb.FriendApplyResp, error) {
	// todo: add your logic here and delete this line

	return &pb.FriendApplyResp{}, nil
}
